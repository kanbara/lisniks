package ui

import (
	"fmt"
	"github.com/kanbara/lisniks/pkg/word"
	"regexp"
)

type ansicolour int

// can't be assed
const light = 10

// TODO could use inverted too
const (
	ANSIBold   int = 1
	ANSIUnder  int = 4
	ANSIInvert int = 7
)

const (
	Black, DarkGrey ansicolour = iota, iota + light
	Red, LightRed
	Green, LightGreen
	Brown, Yellow
	Blue, LightBlue
	Magenta, LightMagenta
	Cyan, LightCyan
	LightGrey, White
)

const (
	beginner       string = "\033[3%dm"
	beginnerStyled string = "\033[3%d;%dm"
	terminator     string = "\033[0m"
)

func StripANSI(str string) string {
	re := regexp.MustCompile(`\033\[[[:digit:]]+;?[[:digit:]]*m`)
	return string(re.ReplaceAll([]byte(str), []byte("")))
}

func ApplyColour(c ansicolour, str string, style int) string {

	switch style {
	case ANSIBold:
		b := fmt.Sprintf(beginner, c)
		if c > LightGrey { // last colour value before brights
			c -= light
			b = fmt.Sprintf(beginnerStyled, c, style)
		}
		return b + str + terminator
	case ANSIInvert, ANSIUnder:
		b := fmt.Sprintf(beginnerStyled, c, style)
		return b + str + terminator
	}

	return str
}

func staticWordClassMap() map[word.Class]ansicolour {
	return map[word.Class]ansicolour{
		// gender
		word.Class{Class: 2, Value: 0}: Magenta,   // masc
		word.Class{Class: 2, Value: 2}: Green,     // fem
		word.Class{Class: 2, Value: 4}: LightGrey, // neut
		// declension
		word.Class{Class: 3, Value: 0}:  Blue,         // a-stem
		word.Class{Class: 3, Value: 16}: Green,        // r-stem
		word.Class{Class: 3, Value: 2}:  Red,          // o-stem
		word.Class{Class: 3, Value: 18}: LightMagenta, // ja-stem
		word.Class{Class: 3, Value: 19}: LightCyan,    // jo-stem
		word.Class{Class: 3, Value: 4}:  Cyan,         // i-stem
		word.Class{Class: 3, Value: 21}: White,        // irregular
		word.Class{Class: 3, Value: 6}:  Brown,        // u-stem
		word.Class{Class: 3, Value: 8}:  LightRed,     // c-stem
		word.Class{Class: 3, Value: 10}: Yellow,       // an-stem
		word.Class{Class: 3, Value: 12}: LightBlue,    // on-stem
		word.Class{Class: 3, Value: 14}: LightGreen,   // in-stem
		// conjugation
		word.Class{Class: 4, Value: 0}:  Cyan,         // weak-j
		word.Class{Class: 4, Value: 2}:  Brown,        // weak-n
		word.Class{Class: 4, Value: 4}:  Green,        // weak-o
		word.Class{Class: 4, Value: 6}:  Magenta,      // weak-a
		word.Class{Class: 4, Value: 8}:  LightRed,     // strong-1
		word.Class{Class: 4, Value: 10}: LightGreen,   // strong-2
		word.Class{Class: 4, Value: 12}: Yellow,       // strong-3
		word.Class{Class: 4, Value: 14}: LightBlue,    // strong-4
		word.Class{Class: 4, Value: 16}: LightMagenta, // strong-5
		word.Class{Class: 4, Value: 18}: LightCyan,    // strong-6
		word.Class{Class: 4, Value: 20}: DarkGrey,     // strong-7
		word.Class{Class: 4, Value: 22}: LightGrey,    // preterite-present
		word.Class{Class: 4, Value: 23}: White,        // irregular
		// aspect
		word.Class{Class: 5, Value: 0}: Blue,  // perfective
		word.Class{Class: 5, Value: 2}: Brown, // imperfective
	}
}

// TODO could make suffixes underlined
func staticPOSMap() map[int64]ansicolour {
	return map[int64]ansicolour{
		6:  Magenta,      // adj suffix
		5:  Magenta,      // adj
		7:  Green,        // adv
		13: Brown,        // conj
		10: Yellow,       // demonst
		11: LightGrey,    // indef pronoun
		14: LightCyan,    // interj
		9:  LightGrey,    // interr. pronoun
		3:  Cyan,         // nom suffix
		2:  Cyan,         // noun
		34: LightMagenta, // numeral
		35: LightMagenta, // ordinal
		8:  LightGrey,    // personal pronoun
		15: Magenta,      // poss adj
		17: LightGrey,    // prefix
		12: Red,          // prepos
		4:  Blue,         // verb
		16: Blue,         // verbal suffix
	}
}

// TODO find a better spot / way to do or supply this (e.g. could be a file)
func WordGrammarColour(str string, wg word.Class) string {

	colour, ok := staticWordClassMap()[wg]
	if !ok {
		// colour was not found, return normal
		return ApplyColour(LightGrey, str, ANSIBold)
	}

	return ApplyColour(colour, str, ANSIBold)
}

// TODO find a better spot / way to do or supply this (e.g. could be a file)
func POSColour(str string, pos int64) string {

	colour, ok := staticPOSMap()[pos]
	if !ok {
		// colour was not found, return white
		return ApplyColour(LightGrey, str, ANSIUnder)
	}

	return ApplyColour(colour, str, ANSIUnder)
}
