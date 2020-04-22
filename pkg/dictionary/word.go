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

// TODO Also apparently we need to filter out the class that doesn't exist e.g.
// hassa is noun which should have gender and declensino, not CLASS4
// there should be type hints on the TypeID of the word grammar class
// a la ApplyTypes.
//
// yay more map lookups
func (d *Dictionary) HumanReadableWordClasses(classes []word.Class) string {
	var out string
	for _, c := range classes {
		val := d.WordGrammar.Get(c)
		out += fmt.Sprintf("%v\n", val.ValueName)
	}

	return out
}

func (d *Dictionary) PrettyWord(w *word.Word) string {
	return fmt.Sprintf("%v (%v) [%v] #%v\n%v\n\tdef: %v\n",
		w.Con,
		w.Local,
		d.PartsOfSpeech.Get(w.Type),
		w.WordID,
		d.HumanReadableWordClasses(w.Classes),
		w.Def)
}