package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

func (m *Manager) NewWordGrammarView(g* gocui.Gui) error {

	maxX, _ := g.Size()

	if v, err := g.SetView(wordGrammarView, 21, 6, maxX-1, 10, 0); err != nil {
		if !gocui.IsUnknownView(err) {
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