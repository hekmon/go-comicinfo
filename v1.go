package comicinfo

import (
	"errors"
	"fmt"
)

func isV1YesNoValid(value string) bool {
	switch value {
	case "", Unknown, No, Yes:
		return true
	default:
		return false
	}
}

type PageV1 struct {
	Image       int        `xml:"Image,attr"`
	Type        PageTypeV1 `xml:"Type,attr"`
	DoublePage  bool       `xml:"DoublePage,attr"`
	ImageSize   int64      `xml:"ImageSize,attr"`
	Key         string     `xml:"Key,attr"`
	ImageWidth  int        `xml:"ImageWidth,attr"`
	ImageHeight int        `xml:"ImageHeight,attr"`
}

func (p *PageV1) Validate() (err error) {
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

type PageTypeV1 string

const (
	FrontCover    PageTypeV1 = "FrontCover"
	InnerCover    PageTypeV1 = "InnerCover"
	Roundup       PageTypeV1 = "Roundup"
	Story         PageTypeV1 = "Story"
	Advertisement PageTypeV1 = "Advertisement"
	Editorial     PageTypeV1 = "Editorial"
	Letters       PageTypeV1 = "Letters"
	Preview       PageTypeV1 = "Preview"
	BackCover     PageTypeV1 = "BackCover"
	Other         PageTypeV1 = "Other"
	Deleted       PageTypeV1 = "Deleted"
)

func (pt PageTypeV1) Valid() bool {
	switch pt {
	case FrontCover, InnerCover, Roundup, Story, Advertisement, Editorial, Letters, Preview, BackCover, Other, Deleted:
		return true
	default:
		return false
	}
}
