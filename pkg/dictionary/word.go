package dictionary

import (
	"fmt"
	"github.com/kanbara/lisniks/pkg/polyglot/word"
)

func (d *Dictionary) PrettyWordStringByID(id int64) string {
	w := d.Lexicon.GetByID(id)
	return d.PrettyWord(w)
}

func (d *Dictionary) PrettyWordStringByLoc(loc int) string {
	w := d.Lexicon.At(loc)
	return d.PrettyWord(w)
}

type HumanClass struct {
	Name  string
	Class word.Class
}

func (d *Dictionary) HumanReadableWordClasses(wordType int64, classes []word.Class) []HumanClass {
	var out []HumanClass
	for _, c := range classes {
		val := d.WordGrammar.Get(wordType, c)
		if val == nil {
			continue
		}

		out = append(out, HumanClass{
			Name: val.ValueName,
			Class: c,
		})
	}

	return out
}

func (d *Dictionary) PrettyWord(w *word.Word) string {
	return fmt.Sprintf("%v (%v) [%v] #%v\n%v\n\tdef: %v\n",
		w.Austrian,
		w.English,
		d.PartsOfSpeech.GetByID(w.Type),
		w.WordID,
		d.HumanReadableWordClasses(w.Type, w.Classes),
		w.Def)
}
