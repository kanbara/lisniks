package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func (m *Manager) NewWordGrammarView(g* gocui.Gui) error {

	maxX, _ := g.Size()

	if v, err := g.SetView(wordGrammarView, 21, 3, maxX-1, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = wordGrammarView
		err := m.UpdateWordGrammarView(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) UpdateWordGrammarView(v *gocui.View) error {
	v.Clear()
	_, err := fmt.Fprintln(v, m.dict.HumanReadableWordClasses(
		m.state.CurrentWord().Type,
		m.state.CurrentWord().Classes))
	if err != nil {
		return err
	}

	return nil
}