package ui

import (
	"github.com/awesome-gocui/gocui"
)

func (m *Manager) SetKeybindings(g *gocui.Gui) error {
	// todo loop thru keybindings too
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitfn); err != nil {
		return err
	}

	if err := g.SetKeybinding("", '/', gocui.ModNone, toSearchView); err != nil {
		return err
	}

	return nil
}

// TODO move these things someplace better
func quitfn(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func toSearchView(g *gocui.Gui, _ *gocui.View) error {
	g.Cursor = true
	return toView(g, searchView)
}

func toView(g *gocui.Gui, view string) error {
	_, err := g.SetCurrentView(view)
	if err != nil {
		return err
	}

	return nil
}
