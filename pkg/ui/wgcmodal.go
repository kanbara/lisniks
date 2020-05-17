package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/kanbara/lisniks/pkg/polyglot/word"
)

type WordGrammarSelectView struct {
	ListView
}

func (w *WordGrammarSelectView) New(name string) error {
	x, y := w.g.Size()
	if v, err := w.g.SetView(name, (x/2)+5, 5, x-5, y-5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = w.viewName
		v.Frame = true
		v.FgColor = gocui.ColorRed
		v.Highlight = true

		if err := w.Update(v); err != nil {
			return err
		}
	}

	return nil
}

func (w *WordGrammarSelectView) Update(v *gocui.View) error {
	v.Clear()

	pos := w.State.SearchState.POSList[w.State.SearchState.SelectedPOS]

	w.Log.Debugf("updated WGS")
	m := w.Dict.WordGrammar.GetAllByType(pos.ID)
	w.itemLen = func() int { return len(m.Values) }

	w.Log.Debugf("for POS %v got %v and len %v", pos.ID, m, w.itemLen())

	for _, t := range m.Values {
		_, err := fmt.Fprintln(v, WordGrammarColour(t.ValueName,
			word.Class{Class: *m.ClassID, Value: t.ValueID}))

		if err != nil {
			return err
		}
	}

	return nil
}