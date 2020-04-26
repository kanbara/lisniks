package ui

import (
	"github.com/awesome-gocui/gocui"
)

func (m *Manager) SetGlobalKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitfn); err != nil {
		return err
	}

	if err := g.SetKeybinding("", '/', gocui.ModNone, toSearchView); err != nil {
		return err
	}

	return nil
}

// used to send the `error` which quits the program
func quitfn(_ *gocui.Gui, _ *gocui.View) error { return gocui.ErrQuit }
