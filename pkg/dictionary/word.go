package dictionary

import (
	"fmt"
	"github.com/kanbara/lisniks/pkg/word"
)

func (d *Dictionary) PrettyWordStringByID(id int64) string {
	w := d.Lexicon.GetByID(id)
	return d.PrettyWord(w)
}

func (d *Dictionary) PrettyWordStringByLoc(loc int) string {
	w := d.Lexicon.At(loc)
	return d.PrettyWord(w)
}

func (d *Dictionary) HumanReadableWordClasses(wordType int64, classes []word.Class) string {
	var out string
	for _, c := range classes {
		val := d.WordGrammar.Get(wordType, c)
		if val == nil {
			continue
		}

		out += fmt.Sprintf("%v ", val.ValueName)
	}

	return out
}

func (d *Dictionary) PrettyWord(w *word.Word) string {
	return fmt.Sprintf("%v (%v) [%v] #%v\n%v\n\tdef: %v\n",
		w.Con,
		w.Local,
		d.PartsOfSpeech.Get(w.Type),
		w.WordID,
		d.HumanReadableWordClasses(w.Type, w.Classes),
		w.Def)
}
