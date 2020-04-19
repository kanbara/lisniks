package dictionary

import (
	"time"
)

type Dictionary struct {
	Version            string                     `xml:"PolyGlotVer"`
	LastUpdated        time.Time                  `xml:"DictSaveDate"`
	LanguageProperties LanguageProperties         `xml:"languageProperties"`
	WordGrammarClasses WordGrammarClassCollection `xml:"wordGrammarClassCollection"`
	PartsOfSpeech      PartsOfSpeech              `xml:"partsOfSpeech"`
	Lexicon            Lexicon                    `xml:"lexicon>word"`
	Etymologies        EtymologyCollection        `xml:"etymologyCollection"`
	Declensions        DeclensionCollection       `xml:"declensionCollection"`
}

type LanguageProperties struct {
	Name string `xml:"langName"`
	Copyright string `xml:"langPropAuthorCopyright"`
}

// Version gets the data from the Copyright field, as we assume that's where the version is stored
func (l LanguageProperties) Version() string {
	return l.Copyright
}

type WordGrammarClassCollection struct {
}

type PartsOfSpeech struct {
}

type EtymologyCollection struct {
}

type DeclensionCollection struct {
}
