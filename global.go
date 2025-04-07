package comicinfo

import (
	"golang.org/x/text/language"
)

const (
	ComicInfoFileName = "ComicInfo.xml"
	xmlnsxni          = "http://www.w3.org/2001/XMLSchema-instance"
)

var (
	// LanguageEnglish is the standard English language ISO code. Available as a helper/shortcut.
	LanguageEnglish = language.English.String()
)
