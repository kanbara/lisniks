package dictionary

import (
	"time"
)

type Dictionary struct {
	Version            Rawstring            `xml:"PolyGlotVer"`
	LastUpdated        time.Time            `xml:"DictSaveDate"`
	LanguageProperties LanguageProperties   `xml:"languageProperties"`
	WordGrammarClasses WordGrammarClass	 	`xml:"wordGrammarClassCollection>wordGrammarClassNode"`
	PartsOfSpeech      PartsOfSpeech        `xml:"partsOfSpeech>class"`
	Lexicon            Lexicon              `xml:"lexicon>word"`
	Etymologies        EtymologyCollection  `xml:"etymologyCollection"`
	Declensions        DeclensionCollection `xml:"declensionCollection"`
}

type LanguageProperties struct {
	Name Rawstring `xml:"langName"`
	Copyright Rawstring `xml:"langPropAuthorCopyright"`
}

// Version gets the data from the Copyright field, as we assume that's where the version is stored
func (l LanguageProperties) Version() string {
	return l.Copyright.String()
}

type EtymologyCollection struct {
}

type DeclensionCollection struct {
}