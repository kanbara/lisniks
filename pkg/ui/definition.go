package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func (m *Manager) NewDefinitionView(g* gocui.Gui) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView(defnView, 21, 11, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
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
	_, err := fmt.Fprintln(v, m.state.CurrentWord().Def)
	if err != nil {
		return err
	}

	return nil
}