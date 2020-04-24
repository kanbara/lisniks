package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func (m *Manager) NewLexiconView(g *gocui.Gui) error {

	_, maxY := g.Size()
	if v, err := g.SetView(lexView, 0, 3, 20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "lexicon"
		v.Frame = true

		for _, w  := range m.dict.Lexicon.Words() {
			_, err := fmt.Fprintln(v, w.Con)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
