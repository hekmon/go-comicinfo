package comicinfo

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
)

const (
	v2SchemaLocationURL = "https://raw.githubusercontent.com/anansi-project/comicinfo/refs/heads/main/schema/v2.0/ComicInfo.xsd"
)

// ComicInfov2 represents the structure of a version 2 ComicInfo.xml file.
type ComicInfov2 struct {
	ComicInfov1
	Day                 int      `xml:"Day,omitempty"`                 // Usually contains the release date of the book.
	Characters          string   `xml:"Characters,omitempty"`          // Characters present in the book. It is accepted that multiple values are comma separated.
	Teams               string   `xml:"Teams,omitempty"`               // Teams present in the book. Usually refer to super-hero teams (e.g. Avengers). It is accepted that multiple values are comma separated.
	Locations           string   `xml:"Locations,omitempty"`           // Locations mentioned in the book. It is accepted that multiple values are comma separated.
	ScanInformation     string   `xml:"ScanInformation,omitempty"`     // A free text field, usually used to store information about who scanned the book.
	StoryArc            string   `xml:"StoryArc,omitempty"`            // The story arc that books belong to. For example, for Undiscovered Country, issues 1-6 are part of the Destiny story arc, issues 7-12 are part of the Unity story arc.
	SeriesGroup         string   `xml:"SeriesGroup,omitempty"`         // A group or collection the series belongs to. It is accepted that multiple values are comma separated.
	AgeRating           string   `xml:"AgeRating,omitempty"`           // The age rating of the book. Possible values are "Unknown", "Everyone", "Teen", "Mature", "Adults Only 18+", "Not Yet Rated".
	Pages               PagesV2  `xml:"Pages,omitempty"`               // Pages of the comic book. Each page should have an Image element with a file path to the image.
	CommunityRating     *float64 `xml:"CommunityRating,omitempty"`     // Community rating of the book, from 0.0 to 5.0, 2 digits allowed.
	MainCharacterOrTeam string   `xml:"MainCharacterOrTeam,omitempty"` // Main character or team mentioned in the book. It is accepted that a single value should be present.
	Review              string   `xml:"Review,omitempty"`              // Review of the book.
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
func (ci ComicInfov2) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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
	// Validate v1 fields
	if err = ci.ComicInfov1.Validate(); err != nil {
		return err
	}
	if len(ci.ComicInfov1.Pages) > 0 {
		return errors.New("pages v1 should not be set in ComicInfov2")
	}
	// Age Rating
	if !isV2AgeRatingValid(ci.AgeRating) {
		return fmt.Errorf("failed to validate AgeRating: unknown value %q", ci.AgeRating)
	}
	// Pages
	if err = ci.Pages.Validate(); err != nil {
		return fmt.Errorf("failed to validate Pages: %w", err)
	}
	// Community Rating
	if !isV2CommunityRatingValid(ci.CommunityRating) {
		return fmt.Errorf("failed to validate CommunityRating: invalid value %f", *ci.CommunityRating)
	}
	return
}

func isV2MangaValid(value string) bool {
	switch value {
	case "", Unknown, No, Yes, YesAndRightToLeft:
		return true
	default:
		return false
	}
}

func isV2AgeRatingValid(value string) bool {
	switch value {
	case "", Unknown, AgeRatingAdultsOnly18Plus, AgeRatingEarlyChildhood, AgeRatingEveryone, AgeRatingEveryone10Plus,
		AgeRatingG, AgeRatingKidsToAdults, AgeRatingM, AgeRatingMA15Plus, AgeRatingMature17Plus, AgeRatingPG,
		AgeRatingR18Plus, AgeRatingRatingPending, AgeRatingTeen, AgeRatingX18Plus:
		return true
	default:
		return false
	}
}

type PagesV2 []PageV2

func (ps PagesV2) Validate() (err error) {
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

type PageV2 struct {
	PageV1
	Bookmark string `xml:"Bookmark,attr"`
}

func isV2CommunityRatingValid(value *float64) bool {
	if value == nil {
		return true
	}
	if *value < 0 {
		return false
	}
	if *value > 5 {
		return false
	}
	// 2 digits allowed
	if math.Abs(*value-math.Round(*value*100)/100) > 0.01 {
		return false
	}
	return true
}
