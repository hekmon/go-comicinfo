package comicinfo

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
)

const (
	v21SchemaLocationURL = "https://github.com/anansi-project/comicinfo/raw/refs/heads/main/drafts/v2.1/ComicInfo.xsd"
)

// ComicInfov21 represents the structure of a version 2.1 DRAFT ComicInfo.xml file.
type ComicInfov21 struct {
	ComicInfov2
	Translator      string   `xml:"Translator,omitempty"`      // A person or organization who renders a text from one language into another, or from an older form of a language into the modern form. This can also be used for fan translations ("scanlator").
	Tags            string   `xml:"Tags,omitempty"`            // Tags of the book or series. For example, ninja or school life. It is accepted that multiple values are comma separated.
	CommunityRating *float64 `xml:"CommunityRating,omitempty"` // Community rating of the book, from 0.0 to 5.0, 1 digits allowed.
	StoryArcNumber  string   `xml:"StoryArcNumber,omitempty"`  // While StoryArc was originally designed to store the arc within a series, it was often used to indicate that a book was part of a reading order, composed of books from multiple series. Mylar for instance was using the field as such. Since StoryArc itself wasn't able to carry the information about ordering of books within a reading order, StoryArcNumber was added. StoryArc and StoryArcNumber can work in combination, to indicate in which position the book is located at for a specific reading order. It is accepted that multiple values can be specified for both StoryArc and StoryArcNumber. Multiple values are comma separated.
	GTIN            string   `xml:"GTIN,omitempty"`            // A Global Trade Item Number identifying the book. GTIN incorporates other standards like ISBN, ISSN, EAN, or JAN.
	Pages           PagesV2  `xml:"Pages,omitempty"`           // Pages of the comic book. Each page should have an Image element with a file path to the image.
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
	// Validate v2 fields
	if err = ci.ComicInfov2.Validate(); err != nil {
		return err
	}
	// Community Rating
	if !isV21CommunityRatingValid(ci.CommunityRating) {
		return fmt.Errorf("failed to validate CommunityRating: invalid value %f", *ci.CommunityRating)
	}
	return
}

func isV21CommunityRatingValid(value *float64) bool {
	if value == nil {
		return true
	}
	if *value < 0 {
		return false
	}
	if *value > 5 {
		return false
	}
	// 1 digits allowed
	if math.Abs(*value-math.Round(*value*10)/10) > 0.1 {
		return false
	}
	return true
}
