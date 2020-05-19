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
		v.Highlight = true

		// todo we need this rn because we create the view again but don't reset the sel
		w.State.SearchState.SelectedPOS = 0

		if err := w.Update(v); err != nil {
			return err
		}
	}

	return nil
}

func (w *WordGrammarSelectView) Update(v *gocui.View) error {
	v.Clear()

	if w.State.SearchState.SelectedPOS >= len(w.State.SearchState.POSList) {
		w.Log.Debugf("oops got %v", w.State.SearchState.SelectedPOS)
	}

	pos := w.State.SearchState.POSList[w.State.SearchState.SelectedPOS]

	classID, m := w.Dict.WordGrammar.GetAllByType(pos.ID)
	w.itemLen = func() int { return len(m) }

	if w.itemLen() != 0 {
		w.Log.Debugf("for POS %v got %v and len %v", pos.ID, m, w.itemLen())
	}

	for _, t := range m {
		_, err := fmt.Fprintln(v, WordGrammarColour(t.ValueName,
			word.Class{Class: *classID, Value: t.ValueID}))

		if err != nil {
			return err
		}
	}

	return nil
}