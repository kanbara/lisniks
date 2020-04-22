package dictionary

import (
	"github.com/kanbara/lisniks/pkg/lexicon"
	"github.com/kanbara/lisniks/pkg/partsofspeech"
	"github.com/kanbara/lisniks/pkg/strings"
	"github.com/kanbara/lisniks/pkg/wordgrammar"
	"time"
)

type File struct {
	Version            strings.Rawstring           `xml:"PolyGlotVer"`
	LastUpdated        time.Time                   `xml:"DictSaveDate"`
	LanguageProperties LanguageProperties          `xml:"languageProperties"`
	WordGrammarClasses wordgrammar.Class           `xml:"wordGrammarClassCollection>wordGrammarClassNode"`
	PartsOfSpeech      partsofspeech.PartsOfSpeech `xml:"partsOfSpeech>class"`
	Lexicon            lexicon.Lexicon             `xml:"lexicon>word"`
	Etymologies        EtymologyCollection         `xml:"etymologyCollection"`
	Declensions        DeclensionCollection        `xml:"declensionCollection"`
}

type LanguageProperties struct {
	Name      strings.Rawstring `xml:"langName"`
	Copyright strings.Rawstring `xml:"langPropAuthorCopyright"`
}

// Version gets the data from the Copyright field, as we assume that's where the version is stored
func (l LanguageProperties) Version() string {
	return l.Copyright.String()
}

type EtymologyCollection struct {
}

type DeclensionCollection struct {
}