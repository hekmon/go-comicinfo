package comicinfo

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"golang.org/x/text/language"
)

const (
	v1SchemaLocationURL = "https://github.com/anansi-project/comicinfo/raw/refs/heads/main/schema/v1.0/ComicInfo.xsd"
)

// ComicInfoComicInfov1 represents the structure of a version 1 ComicInfo.xml file.
type ComicInfov1 struct {
	Title           string `xml:"Title,omitempty"`           // Title of the book.
	Series          string `xml:"Series,omitempty"`          // Title of the series the book is part of.
	Number          int    `xml:"Number,omitempty"`          // Number of the book in the series.
	Count           int    `xml:"Count,omitempty"`           // The total number of books in the series. The Count could be different on each book in a series. Consuming applications should consider using only the value for the latest book in the series.
	Volume          int    `xml:"Volume,omitempty"`          // Volume containing the book. Volume is a notion that is specific to US Comics, where the same series can have multiple volumes. Volumes can be referenced by number (1, 2, 3…) or by year (2018, 2020…).
	AlternateSeries string `xml:"AlternateSeries,omitempty"` // Quite specific to US comics, some books can be part of cross-over story arcs. Those fields can be used to specify an alternate series, its number and count of books.
	AlternateNumber int    `xml:"AlternateNumber,omitempty"` // Quite specific to US comics, some books can be part of cross-over story arcs. Those fields can be used to specify an alternate series, its number and count of books.
	AlternateCount  int    `xml:"AlternateCount,omitempty"`  // Quite specific to US comics, some books can be part of cross-over story arcs. Those fields can be used to specify an alternate series, its number and count of books.
	Summary         string `xml:"Summary,omitempty"`         // A description or summary of the book.
	Notes           string `xml:"Notes,omitempty"`           // A free text field, usually used to store information about the application that created the ComicInfo.xml file.
	Year            int    `xml:"Year,omitempty"`            // Usually contains the release date of the book.
	Month           int    `xml:"Month,omitempty"`           // Usually contains the release date of the book.
	Writer          string `xml:"Writer,omitempty"`          // Person or organization responsible for creating the scenario. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Penciller       string `xml:"Penciller,omitempty"`       // Person or organization responsible for drawing the art. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Inker           string `xml:"Inker,omitempty"`           // Person or organization responsible for inking the pencil art. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Colorist        string `xml:"Colorist,omitempty"`        // Person or organization responsible for applying color to drawings. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Letterer        string `xml:"Letterer,omitempty"`        // Person or organization responsible for drawing text and speech bubbles. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	CoverArtist     string `xml:"CoverArtist,omitempty"`     // Person or organization responsible for drawing the cover art. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Editor          string `xml:"Editor,omitempty"`          // A person or organization contributing to a resource by revising or elucidating the content, e.g., adding an introduction, notes, or other critical matter. An editor may also prepare a resource for production, publication, or distribution. In order to cater for multiple creator with the same role, it is accepted that values are comma separated.
	Publisher       string `xml:"Publisher,omitempty"`       // A person or organization responsible for publishing, releasing, or issuing a resource.
	Imprint         string `xml:"Imprint,omitempty"`         // An imprint is a group of publications under the umbrella of a larger imprint or a Publisher. For example, Vertigo is an Imprint of DC Comics.
	Genre           string `xml:"Genre,omitempty"`           // Genre of the book or series. For example, Science-Fiction or Shonen. It is accepted that multiple values are comma separated.
	Web             string `xml:"Web,omitempty"`             // A URL pointing to a reference website for the book. It is accepted that multiple values are space separated (as spaces in URL will be encoded as %20).
	PageCount       int    `xml:"PageCount,omitempty"`       // The number of pages in the book.
	Language        string `xml:"LanguageISO,omitempty"`     // ISO code of the language the book is written in. You can use "golang.org/x/text/language" to get valid codes, eg language.English.String()
	Format          string `xml:"format,omitempty"`          // The original publication's binding format for scanned physical books or presentation format for digital sources. "TBP", "HC", "Web", "Digital" are common designators.
	BlackAndWhite   YesNo  `xml:"BlackAndWhite,omitempty"`   // Whether the book is in black and white.
	Manga           Manga  `xml:"Manga,omitempty"`           // Whether the book is a manga. This also defines the reading direction as right-to-left when set to YesAndRightToLeft.
	Pages           Pages  `xml:"Pages,omitempty"`           // Pages of the comic book. Each page should have an Image element with a file path to the image.
}

// Encode will produce a ComicInfo v2 XML content. It will validate the ComicInfo struct before encoding it into XML format.
func (ci ComicInfov1) Encode(output io.Writer) (err error) {
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
		return fmt.Errorf("failed to encode ComicInfo v1 XML: %w", err)
	}
	return
}

// MarshalXML implements the xml.Marshaler interface to automatically add schema attributes.
// User should use Encode() instead of this method directly. This method is used internally by Encode().
func (ci ComicInfov1) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type Mask ComicInfov1
	type attr struct {
		Mask
		XSI            string `xml:"xmlns:xsi,attr"`
		SchemaLocation string `xml:"xsi:schemaLocation,attr"`
	}
	return e.EncodeElement(attr{
		Mask:           Mask(ci),
		XSI:            xmlnsxni,
		SchemaLocation: v1SchemaLocationURL,
	}, start)
}

// Validate checks if some of the fields with particular constraints are valid. It returns an error if any field fails validation.
func (ci ComicInfov1) Validate() (err error) {
	// URL(s)
	for index, URL := range strings.Split(ci.Web, " ") {
		if _, err = url.Parse(URL); err != nil {
			return fmt.Errorf("failed to validate URL #%d: %w", index, err)
		}
	}
	// Language
	if ci.Language != "" {
		if _, err = language.Parse(ci.Language); err != nil {
			return fmt.Errorf("failed to validate Language: %s", ci.Language)
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
	// Pages
	if err = ci.Pages.Validate(); err != nil {
		return fmt.Errorf("failed to validate Pages: %w", err)
	}
	return
}

type YesNo string

const (
	Unknown YesNo = "Unknown"
	No      YesNo = "No"
	Yes     YesNo = "Yes"
)

func (yn YesNo) IsValid() bool {
	switch yn {
	case "", Unknown, No, Yes:
		return true
	default:
		return false
	}
}

type Manga string

const (
	MangaUnknown           Manga = "Unknown"
	MangaNo                Manga = "No"
	MangaYes               Manga = "Yes"
	MangaYesAndRightToLeft Manga = "YesAndRightToLeft"
)

func (m Manga) IsValid() bool {
	switch m {
	case "", MangaUnknown, MangaNo, MangaYes, MangaYesAndRightToLeft:
		return true
	default:
		return false
	}
}

type Pages []Page

func (ps Pages) Validate() (err error) {
	keys := make(map[string]struct{}, len(ps))
	var ok bool
	for i, p := range ps {
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

type Page struct {
	Image       int      `xml:"Image,attr"`
	Type        PageType `xml:"Type,attr"`
	DoublePage  bool     `xml:"DoublePage,attr"`
	ImageSize   int      `xml:"ImageSize,attr"`
	Key         string   `xml:"Key,attr"`
	ImageWidth  int      `xml:"ImageWidth,attr"`
	ImageHeight int      `xml:"ImageHeight,attr"`
}

func (p *Page) Validate() (err error) {
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

type PageType string

const (
	PageTypeFrontCover    PageType = "FrontCover"
	PageTypeInnerCover    PageType = "InnerCover"
	PageTypeRoundup       PageType = "Roundup"
	PageTypeStory         PageType = "Story"
	PageTypeAdvertisement PageType = "Advertisement"
	PageTypeEditorial     PageType = "Editorial"
	PageTypeLetters       PageType = "Letters"
	PageTypePreview       PageType = "Preview"
	PageTypeBackCover     PageType = "BackCover"
	PageTypeOther         PageType = "Other"
	PageTypeDeleted       PageType = "Deleted"
)

func (pt PageType) Valid() bool {
	switch pt {
	case PageTypeFrontCover, PageTypeInnerCover, PageTypeRoundup, PageTypeStory, PageTypeAdvertisement,
		PageTypeEditorial, PageTypeLetters, PageTypePreview, PageTypeBackCover, PageTypeOther, PageTypeDeleted:
		return true
	default:
		return false
	}
}
