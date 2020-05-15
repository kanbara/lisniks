package dictionary

import (
	"github.com/kanbara/lisniks/pkg/polyglot/word"
)

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