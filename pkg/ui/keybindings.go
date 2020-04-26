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

	if err := g.SetKeybinding(lexView, gocui.KeyArrowDown, gocui.ModNone, m.nextWord); err != nil {
		return err
	}

	if err := g.SetKeybinding(lexView, gocui.KeyArrowUp, gocui.ModNone, m.prevWord); err != nil {
		return err
	}

	if err := g.SetKeybinding(lexView, gocui.KeyCtrlF, gocui.ModNone, m.nextWordJump); err != nil {
		return err
	}

	if err := g.SetKeybinding(lexView, gocui.KeyCtrlB, gocui.ModNone, m.prevWordJump); err != nil {
		return err
	}

	if err := g.SetKeybinding(searchView, gocui.KeyEnter, gocui.ModNone, m.execSearch); err != nil {
		return err
	}

	if err := g.SetKeybinding(searchView, gocui.KeyEsc, gocui.ModNone, cancelToLexView); err != nil {
		return err
	}

	return nil
}

// TODO move these things someplace better
func quitfn(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func cancelToLexView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	v.Clear()
	err := v.SetCursor(0, 0)
	if err != nil {
		return err
	}

	return toView(g, lexView)
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
