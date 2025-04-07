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
	v21SchemaLocationURL = "https://github.com/anansi-project/comicinfo/raw/refs/heads/main/drafts/v2.1/ComicInfo.xsd"
)

// ComicInfov21 represents the structure of a version 2.1 DRAFT ComicInfo.xml file.
type ComicInfov21 struct {
	Title               string              `xml:"Title,omitempty"`               // Title of the book.
	Series              string              `xml:"Series,omitempty"`              // Title of the series the book is part of.
	Number              int                 `xml:"Number,omitempty"`              // Number of the book in the series.
	Count               int                 `xml:"Count,omitempty"`               // The total number of books in the series. The Count could be different on each book in a series. Consuming applications should consider using only the value for the latest book in the series.
	Volume              int                 `xml:"Volume,omitempty"`              // Volume containing the book. Volume is a notion that is specific to US Comics, where the same series can have multiple volumes. Volumes can be referenced by number (1, 2, 3…) or by year (2018, 2020…).
	AlternateSeries     string              `xml:"AlternateSeries,omitempty"`     // Quite specific to US comics, some books can be part of cross-over story arcs. Those fields can be used to specify an alternate series, its number and count of books.
	AlternateNumber     int                 `xml:"AlternateNumber,omitempty"`     // Quite specific to US comics, some books can be part of cross-over story arcs. Those fields can be used to specify an alternate series, its number and count of books.
	AlternateCount      int                 `xml:"AlternateCount,omitempty"`      // Quite specific to US comics, some books can be part of cross-over story arcs. Those fields can be used to specify an alternate series, its number and count of books.
	Summary             string              `xml:"Summary,omitempty"`             // A description or summary of the book.
	Notes               string              `xml:"Notes,omitempty"`               // A free text field, usually used to store information about the application that created the ComicInfo.xml file.
	Year                int                 `xml:"Year,omitempty"`                // Usually contains the release date of the book.
	Month               int                 `xml:"Month,omitempty"`               // Usually contains the release date of the book.
	Day                 int                 `xml:"Day,omitempty"`                 // Usually contains the release date of the book.
	Publisher           string              `xml:"Publisher,omitempty"`           // A person or organization responsible for publishing, releasing, or issuing a resource.
	Imprint             string              `xml:"Imprint,omitempty"`             // An imprint is a group of publications under the umbrella of a larger imprint or a Publisher. For example, Vertigo is an Imprint of DC Comics.
	Genre               string              `xml:"Genre,omitempty"`               // Genre of the book or series. For example, Science-Fiction or Shonen. It is accepted that multiple values are comma separated.
	Tags                string              `xml:"Tags,omitempty"`                // Tags of the book or series. For example, ninja or school life. It is accepted that multiple values are comma separated.
	Web                 string              `xml:"Web,omitempty"`                 // A URL pointing to a reference website for the book. It is accepted that multiple values are space separated (as spaces in URL will be encoded as %20).
	PageCount           int                 `xml:"PageCount,omitempty"`           // The number of pages in the book.
	LanguageISO         string              `xml:"LanguageISO,omitempty"`         // ISO code of the language the book is written in. You can use "golang.org/x/text/language" to get valid codes, eg language.English.String()
	Format              string              `xml:"format,omitempty"`              // The original publication's binding format for scanned physical books or presentation format for digital sources. "TBP", "HC", "Web", "Digital" are common designators.
	BlackAndWhite       YesNo               `xml:"BlackAndWhite,omitempty"`       // Whether the book is in black and white.
	Manga               Manga               `xml:"Manga,omitempty"`               // Whether the book is a manga. This also defines the reading direction as right-to-left when set to YesAndRightToLeft.
	Characters          string              `xml:"Characters,omitempty"`          // Characters present in the book. It is accepted that multiple values are comma separated.
	Teams               string              `xml:"Teams,omitempty"`               // Teams present in the book. Usually refer to super-hero teams (e.g. Avengers). It is accepted that multiple values are comma separated.
	Locations           string              `xml:"Locations,omitempty"`           // Locations mentioned in the book. It is accepted that multiple values are comma separated.
	ScanInformation     string              `xml:"ScanInformation,omitempty"`     // A free text field, usually used to store information about who scanned the book.
	StoryArc            string              `xml:"StoryArc,omitempty"`            // The story arc that books belong to. For example, for Undiscovered Country, issues 1-6 are part of the Destiny story arc, issues 7-12 are part of the Unity story arc.
	StoryArcNumber      string              `xml:"StoryArcNumber,omitempty"`      // While StoryArc was originally designed to store the arc within a series, it was often used to indicate that a book was part of a reading order, composed of books from multiple series. Mylar for instance was using the field as such. Since StoryArc itself wasn't able to carry the information about ordering of books within a reading order, StoryArcNumber was added. StoryArc and StoryArcNumber can work in combination, to indicate in which position the book is located at for a specific reading order. It is accepted that multiple values can be specified for both StoryArc and StoryArcNumber. Multiple values are comma separated.
	SeriesGroup         string              `xml:"SeriesGroup,omitempty"`         // A group or collection the series belongs to. It is accepted that multiple values are comma separated.
	AgeRating           AgeRating           `xml:"AgeRating,omitempty"`           // The age rating of the book. Possible values are "Unknown", "Everyone", "Teen", "Mature", "Adults Only 18+", "Not Yet Rated".
	Pages               PagesV2             `xml:"Pages,omitempty"`               // Pages of the comic book. Each page should have an Image element with a file path to the image.
	CommunityRating     *CommunityRatingV21 `xml:"CommunityRating,omitempty"`     // Community rating of the book, from 0.0 to 5.0, 2 digits allowed.
	MainCharacterOrTeam string              `xml:"MainCharacterOrTeam,omitempty"` // Main character or team mentioned in the book. It is accepted that a single value should be present.
	Review              string              `xml:"Review,omitempty"`              // Review of the book.
	GTIN                string              `xml:"GTIN,omitempty"`                // A Global Trade Item Number identifying the book. GTIN incorporates other standards like ISBN, ISSN, EAN, or JAN.
	// According to the schema, each creator element can only be present once. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Writer      string `xml:"Writer,omitempty"`      // Person or organization responsible for creating the scenario.
	Penciller   string `xml:"Penciller,omitempty"`   // Person or organization responsible for drawing the art.
	Inker       string `xml:"Inker,omitempty"`       // Person or organization responsible for inking the pencil art.
	Colorist    string `xml:"Colorist,omitempty"`    // Person or organization responsible for applying color to drawings.
	Letterer    string `xml:"Letterer,omitempty"`    // Person or organization responsible for drawing text and speech bubbles.
	CoverArtist string `xml:"CoverArtist,omitempty"` // Person or organization responsible for drawing the cover art.
	Editor      string `xml:"Editor,omitempty"`      // A person or organization contributing to a resource by revising or elucidating the content, e.g., adding an introduction, notes, or other critical matter. An editor may also prepare a resource for production, publication, or distribution.
	Translator  string `xml:"Translator,omitempty"`  // A person or organization who renders a text from one language into another, or from an older form of a language into the modern form. This can also be used for fan translations ("scanlator").
}

// Encode will produce a ComicInfo v2.1 DRAFT XML content. It will validate the ComicInfo struct before encoding it into XML format.
func (ci ComicInfov21) Encode(output io.Writer) (err error) {
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
func (ci ComicInfov21) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Mask ComicInfov21
	type attr struct {
		Mask
		XSI            string `xml:"xmlns:xsi,attr"`
		SchemaLocation string `xml:"xsi:schemaLocation,attr"`
	}
	return e.EncodeElement(attr{
		Mask:           Mask(ci),
		XSI:            xmlnsxni,
		SchemaLocation: v21SchemaLocationURL,
	}, start)
}

// Validate checks if some of the fields with particular constraints are valid. It returns an error if any field fails validation.
func (ci ComicInfov21) Validate() (err error) {
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

type CommunityRatingV21 float64

func (cr *CommunityRatingV21) IsValid() bool {
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
	if math.Abs(float64(*cr)-math.Round(float64(*cr)*10)/10) > 0.1 {
		return false
	}
	return true
}
