package comicinfo

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"strings"

	"golang.org/x/text/language"
)

const (
	v2SchemaLocationURL = "https://raw.githubusercontent.com/anansi-project/comicinfo/refs/heads/main/schema/v2.0/ComicInfo.xsd"
)

// ComicInfov2 represents the structure of a version 2 ComicInfo.xml file.
type ComicInfov2 struct {
	Title               string           `xml:"Title,omitempty"`               // Title of the book.
	Series              string           `xml:"Series,omitempty"`              // Title of the series the book is part of.
	Number              int              `xml:"Number,omitempty"`              // Number of the book in the series.
	Count               int              `xml:"Count,omitempty"`               // The total number of books in the series. The Count could be different on each book in a series. Consuming applications should consider using only the value for the latest book in the series.
	Volume              int              `xml:"Volume,omitempty"`              // Volume containing the book. Volume is a notion that is specific to US Comics, where the same series can have multiple volumes. Volumes can be referenced by number (1, 2, 3…) or by year (2018, 2020…).
	AlternateSeries     string           `xml:"AlternateSeries,omitempty"`     // Quite specific to US comics, some books can be part of cross-over story arcs. Those fields can be used to specify an alternate series, its number and count of books.
	AlternateNumber     int              `xml:"AlternateNumber,omitempty"`     // Quite specific to US comics, some books can be part of cross-over story arcs. Those fields can be used to specify an alternate series, its number and count of books.
	AlternateCount      int              `xml:"AlternateCount,omitempty"`      // Quite specific to US comics, some books can be part of cross-over story arcs. Those fields can be used to specify an alternate series, its number and count of books.
	Summary             string           `xml:"Summary,omitempty"`             // A description or summary of the book.
	Notes               string           `xml:"Notes,omitempty"`               // A free text field, usually used to store information about the application that created the ComicInfo.xml file.
	Year                int              `xml:"Year,omitempty"`                // Usually contains the release date of the book.
	Month               int              `xml:"Month,omitempty"`               // Usually contains the release date of the book.
	Day                 int              `xml:"Day,omitempty"`                 // Usually contains the release date of the book.
	Writer              string           `xml:"Writer,omitempty"`              // Person or organization responsible for creating the scenario. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Penciller           string           `xml:"Penciller,omitempty"`           // Person or organization responsible for drawing the art. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Inker               string           `xml:"Inker,omitempty"`               // Person or organization responsible for inking the pencil art. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Colorist            string           `xml:"Colorist,omitempty"`            // Person or organization responsible for applying color to drawings. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Letterer            string           `xml:"Letterer,omitempty"`            // Person or organization responsible for drawing text and speech bubbles. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	CoverArtist         string           `xml:"CoverArtist,omitempty"`         // Person or organization responsible for drawing the cover art. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Editor              string           `xml:"Editor,omitempty"`              // A person or organization contributing to a resource by revising or elucidating the content, e.g., adding an introduction, notes, or other critical matter. An editor may also prepare a resource for production, publication, or distribution. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Publisher           string           `xml:"Publisher,omitempty"`           // A person or organization responsible for publishing, releasing, or issuing a resource.
	Imprint             string           `xml:"Imprint,omitempty"`             // An imprint is a group of publications under the umbrella of a larger imprint or a Publisher. For example, Vertigo is an Imprint of DC Comics.
	Genre               string           `xml:"Genre,omitempty"`               // Genre of the book or series. For example, Science-Fiction or Shonen. It is accepted that multiple values are comma separated.
	Web                 string           `xml:"Web,omitempty"`                 // A URL pointing to a reference website for the book. It is accepted that multiple values are space separated (as spaces in URL will be encoded as %20).
	PageCount           int              `xml:"PageCount,omitempty"`           // The number of pages in the book.
	LanguageISO         string           `xml:"LanguageISO,omitempty"`         // ISO code of the language the book is written in. You can use "golang.org/x/text/language" to get valid codes, eg language.English.String()
	Format              string           `xml:"Format,omitempty"`              // The original publication's binding format for scanned physical books or presentation format for digital sources. "TBP", "HC", "Web", "Digital" are common designators.
	BlackAndWhite       YesNo            `xml:"BlackAndWhite,omitempty"`       // Whether the book is in black and white.
	Manga               Manga            `xml:"Manga,omitempty"`               // Whether the book is a manga. This also defines the reading direction as right-to-left when set to YesAndRightToLeft.
	Characters          string           `xml:"Characters,omitempty"`          // Characters present in the book. It is accepted that multiple values are comma separated.
	Teams               string           `xml:"Teams,omitempty"`               // Teams present in the book. Usually refer to super-hero teams (e.g. Avengers). It is accepted that multiple values are comma separated.
	Locations           string           `xml:"Locations,omitempty"`           // Locations mentioned in the book. It is accepted that multiple values are comma separated.
	ScanInformation     string           `xml:"ScanInformation,omitempty"`     // A free text field, usually used to store information about who scanned the book.
	StoryArc            string           `xml:"StoryArc,omitempty"`            // The story arc that books belong to. For example, for Undiscovered Country, issues 1-6 are part of the Destiny story arc, issues 7-12 are part of the Unity story arc.
	SeriesGroup         string           `xml:"SeriesGroup,omitempty"`         // A group or collection the series belongs to. It is accepted that multiple values are comma separated.
	AgeRating           AgeRating        `xml:"AgeRating,omitempty"`           // The age rating of the book. Possible values are "Unknown", "Everyone", "Teen", "Mature", "Adults Only 18+", "Not Yet Rated".
	Pages               PagesV2          `xml:"Pages,omitempty"`               // Pages of the comic book. Each page should have an Image element with a file path to the image.
	CommunityRating     *CommunityRating `xml:"CommunityRating,omitempty"`     // Community rating of the book, from 0.0 to 5.0, 2 digits allowed.
	MainCharacterOrTeam string           `xml:"MainCharacterOrTeam,omitempty"` // Main character or team mentioned in the book. It is accepted that a single value should be present.
	Review              string           `xml:"Review,omitempty"`              // Review of the book.
}

// Encode will produce a ComicInfo v2 XML content. It will validate the ComicInfo struct before encoding it into XML format.
func (ci ComicInfov2) Encode(output io.Writer) (err error) {
	if output == nil {
		return errors.New("output cannot be nil")
	}
	// Validate some fields before encoding
	if err = ci.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	// Write header
	if _, err = output.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("failed to write XML header: %w", err)
	}
	// Encode
	encoder := xml.NewEncoder(output)
	encoder.Indent("", "\t")
	if err := encoder.Encode(ci); err != nil {
		return fmt.Errorf("failed to encode ComicInfo v2 XML: %w", err)
	}
	return
}

// MarshalXML implements the xml.Marshaler interface to automatically add schema attributes.
// User should use Encode() instead of this method directly. This method is used internally by Encode().
func (ci ComicInfov2) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "ComicInfo" // Correct name for root name
	type Mask ComicInfov2
	type attr struct {
		Mask
		XSI            string `xml:"xmlns:xsi,attr"`
		SchemaLocation string `xml:"xsi:schemaLocation,attr"`
	}
	return e.EncodeElement(attr{
		Mask:           Mask(ci),
		XSI:            xmlnsxni,
		SchemaLocation: v2SchemaLocationURL,
	}, start)
}

// Validate checks if some of the fields with particular constraints are valid. It returns an error if any field fails validation.
func (ci ComicInfov2) Validate() (err error) {
	// URL(s)
	for index, URL := range strings.Split(ci.Web, " ") {
		if _, err = url.Parse(URL); err != nil {
			return fmt.Errorf("failed to validate URL #%d: %w", index, err)
		}
	}
	// Language
	if ci.LanguageISO != "" {
		if _, err = language.Parse(ci.LanguageISO); err != nil {
			return fmt.Errorf("failed to validate Language: %s", ci.LanguageISO)
		}
	}
	// BlackAndWhite
	if !ci.BlackAndWhite.IsValid() {
		return fmt.Errorf("failed to validate BlackAndWhite: unknown value %q", ci.BlackAndWhite)
	}
	// Manga
	if !ci.Manga.IsValid() {
		return fmt.Errorf("failed to validate Manga: unknown value %q", ci.Manga)
	}
	// Age Rating
	if !ci.AgeRating.IsValid() {
		return fmt.Errorf("failed to validate AgeRating: unknown value %q", ci.AgeRating)
	}
	// Pages
	if err = ci.Pages.Validate(); err != nil {
		return fmt.Errorf("failed to validate Pages: %w", err)
	}
	// Community Rating
	if !ci.CommunityRating.IsValid() {
		return fmt.Errorf("failed to validate CommunityRating: invalid value %f", *ci.CommunityRating)
	}
	return
}

type AgeRating string

const (
	AgeRatingUnknown          AgeRating = "Unknown"
	AgeRatingAdultsOnly18Plus AgeRating = "Adults Only 18+"
	AgeRatingEarlyChildhood   AgeRating = "Early Childhood"
	AgeRatingEveryone         AgeRating = "Everyone"
	AgeRatingEveryone10Plus   AgeRating = "Everyone 10+"
	AgeRatingG                AgeRating = "G"
	AgeRatingKidsToAdults     AgeRating = "Kids to Adults"
	AgeRatingM                AgeRating = "M"
	AgeRatingMA15Plus         AgeRating = "MA15+"
	AgeRatingMature17Plus     AgeRating = "Mature 17+"
	AgeRatingPG               AgeRating = "PG"
	AgeRatingR18Plus          AgeRating = "R18+"
	AgeRatingRatingPending    AgeRating = "Rating Pending"
	AgeRatingTeen             AgeRating = "Teen"
	AgeRatingX18Plus          AgeRating = "X18+"
)

func (ag AgeRating) IsValid() bool {
	switch ag {
	case "", AgeRatingUnknown, AgeRatingAdultsOnly18Plus, AgeRatingEarlyChildhood, AgeRatingEveryone,
		AgeRatingEveryone10Plus, AgeRatingG, AgeRatingKidsToAdults, AgeRatingM, AgeRatingMA15Plus,
		AgeRatingMature17Plus, AgeRatingPG, AgeRatingR18Plus, AgeRatingRatingPending, AgeRatingTeen,
		AgeRatingX18Plus:
		return true
	default:
		return false
	}
}

type PagesV2 struct {
	Pages []PageV2 `xml:"Page"`
}

func (ps PagesV2) Validate() (err error) {
	keys := make(map[string]struct{}, len(ps.Pages))
	var ok bool
	for i, p := range ps.Pages {
		if _, ok = keys[p.Key]; ok {
			return fmt.Errorf("duplicate key found for page %d: %q", i+1, p.Key)
		}
		keys[p.Key] = struct{}{}
		if err = p.Validate(); err != nil {
			return fmt.Errorf("failed to validate page %d: %w", i+1, err)
		}
	}
	return
}

type PageV2 struct {
	Image       int      `xml:"Image,attr"`
	Type        PageType `xml:"Type,attr"`
	DoublePage  bool     `xml:"DoublePage,attr"`
	ImageSize   int      `xml:"ImageSize,attr"`
	Key         string   `xml:"Key,attr"`
	Bookmark    string   `xml:"Bookmark,attr"`
	ImageWidth  int      `xml:"ImageWidth,attr"`
	ImageHeight int      `xml:"ImageHeight,attr"`
}

func (p *PageV2) Validate() (err error) {
	if !p.Type.Valid() {
		return fmt.Errorf("invalid page type: %q", p.Type)
	}
	if !(p.ImageWidth > 0 || p.ImageWidth == -1) {
		return errors.New("image width must be greater than 0 or -1")
	}
	if !(p.ImageHeight > 0 || p.ImageHeight == -1) {
		return errors.New("image height must be greater than 0 or -1")
	}
	return
}

type CommunityRating float64

func (cr *CommunityRating) IsValid() bool {
	if cr == nil {
		return true
	}
	if *cr < 0 {
		return false
	}
	if *cr > 5 {
		return false
	}
	// 2 digits allowed
	if math.Abs(float64(*cr)-math.Round(float64(*cr)*100)/100) > 0.01 {
		return false
	}
	return true
}
