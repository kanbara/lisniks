package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type WordGrammarView DefaultView

func (w *WordGrammarView) New(g *gocui.Gui, name string) error {

	maxX, _ := g.Size()

	if v, err := g.SetView(name, 21, 6, maxX-1, 10, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = name
		err := w.Update(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *WordGrammarView) Update(v *gocui.View) error {
	v.Clear()

	if w.state.CurrentWord() != nil {
		_, err := fmt.Fprintln(v, w.dict.HumanReadableWordClasses(
			w.state.CurrentWord().Type,
			w.state.CurrentWord().Classes))
		if err != nil {
			return err
		}
	}

	return nil
}