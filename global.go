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

const (
	Unknown           = "Unknown"
	No                = "No"
	Yes               = "Yes"
	YesAndRightToLeft = "YesAndRightToLeft"
)

const (
	AgeRatingAdultsOnly18Plus = "Adults Only 18+"
	AgeRatingEarlyChildhood   = "Early Childhood"
	AgeRatingEveryone         = "Everyone"
	AgeRatingEveryone10Plus   = "Everyone 10+"
	AgeRatingG                = "G"
	AgeRatingKidsToAdults     = "Kids to Adults"
	AgeRatingM                = "M"
	AgeRatingMA15Plus         = "MA15+"
	AgeRatingMature17Plus     = "Mature 17+"
	AgeRatingPG               = "PG"
	AgeRatingR18Plus          = "R18+"
	AgeRatingRatingPending    = "Rating Pending"
	AgeRatingTeen             = "Teen"
	AgeRatingX18Plus          = "X18+"
)

