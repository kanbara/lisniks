package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

func (m *Manager) NewCurrentWordView(g *gocui.Gui) error {
	if v, err := g.SetView(currentWordView, 0, 3, 20, 5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Frame = true
		v.FgColor = gocui.ColorGreen

		err := m.updateCurrentWordView(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) updateCurrentWordView(v *gocui.View) error {
	v.Clear()

	if m.state.CurrentWord() != nil {
		v.FgColor = gocui.ColorGreen
		_, err := fmt.Fprintln(v, m.state.CurrentWord().Con)
		if err != nil {
			return err
		}
	} else {
		v.FgColor = gocui.ColorRed
		_, err := fmt.Fprintln(v, "NO WORDS FOUND")
		if err != nil {
			return err
		}
	}

	return nil
}
