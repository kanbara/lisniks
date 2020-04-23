package dictionary

import (
	"github.com/kanbara/lisniks/pkg/declension"
	"github.com/kanbara/lisniks/pkg/language"
	"github.com/kanbara/lisniks/pkg/lexicon"
	"github.com/kanbara/lisniks/pkg/partsofspeech"
	"github.com/kanbara/lisniks/pkg/strings"
	"github.com/kanbara/lisniks/pkg/wordgrammar"
	"time"
)

type File struct {
	Version            strings.Rawstring           `xml:"PolyGlotVer"`
	LastUpdated        time.Time                   `xml:"DictSaveDate"`
	LanguageProperties language.Properties         `xml:"languageProperties"`
	WordGrammarClasses wordgrammar.Class           `xml:"wordGrammarClassCollection>wordGrammarClassNode"`
	PartsOfSpeech      partsofspeech.PartsOfSpeech `xml:"partsOfSpeech>class"`
	Lexicon            lexicon.Lexicon             `xml:"lexicon>word"`
	Etymologies        EtymologyCollection         `xml:"etymologyCollection"`
	Declensions        declension.Declensions      `xml:"declensionCollection>declensionNode"`
}

type EtymologyCollection struct {
}
