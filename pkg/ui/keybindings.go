package ui

import (
	"github.com/jroimartin/gocui"
)

func (m *Manager) SetKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitfn); err != nil {
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

	return nil
}


func quitfn(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}