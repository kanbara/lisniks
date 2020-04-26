package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

func (m *Manager) NewDefinitionView(g* gocui.Gui) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView(defnView, 21, 14, maxX-1, maxY-4, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = defnView
		v.Wrap = true
		err := m.UpdateDefinition(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) UpdateDefinition(v *gocui.View) error {
	v.Clear()

	if m.state.CurrentWord() != nil {
		_, err := fmt.Fprintln(v, m.state.CurrentWord().Def)
		if err != nil {
			return err
		}
	}

	return nil
}