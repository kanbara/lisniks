package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

func (m *Manager) NewLocalWordView(g* gocui.Gui) error {

	maxX, _ := g.Size()

	if v, err := g.SetView(localWordView, 21, 3, maxX-1, 5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = localWordView
		v.FgColor = gocui.ColorYellow
		err := m.UpdateLocalWordView(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) UpdateLocalWordView(v *gocui.View) error {
	v.Clear()
	_, err := fmt.Fprintln(v, m.state.CurrentWord().Local)
	if err != nil {
		return err
	}

	return nil
}