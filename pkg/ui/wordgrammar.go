package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type WordGrammarView struct {
	DefaultView
}

func (w *WordGrammarView) New(name string) error {

	maxX, _ := w.g.Size()

	if v, err := w.g.SetView(name, 21, 6, maxX-1, 9, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Frame = true

		if _, err := w.g.SetViewOnBottom(name); err != nil {
			return err
		}

		err := w.Update(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *WordGrammarView) Update(v *gocui.View) error {
	v.Clear()

	if w.State.CurrentWord() != nil {

		var str string

		// get all applicable human names and classes for this word
		classes := w.Dict.HumanReadableWordClasses(
			w.State.CurrentWord().Type,
			w.State.CurrentWord().Classes)

		for i := 0; i < len(classes); i++ {
			str += WordGrammarColour(classes[i].Name, classes[i].Class)
			if i != len(classes) {
				str += " "
			}
		}

		_, err := fmt.Fprintln(v, str)
		if err != nil {
			return err
		}
	}

	return nil
}
