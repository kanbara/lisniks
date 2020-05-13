package ui

import (
	"github.com/awesome-gocui/gocui"
)

func GlobalKeybindingKeys() []gocui.Key {
	// XXX that's dumb
	// the black/whitelist code only checks the KEY but for runes
	// the KEY is always 0, had to scour the source for this. whatever
	return []gocui.Key{gocui.KeyCtrlR, gocui.KeyCtrlC, 0}
}

func (m *Manager) SetGlobalKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("",
		gocui.KeyCtrlC,
		gocui.ModNone,
		m.quitmodal); err != nil {
		return err
	}

	if err := g.SetKeybinding("",
		gocui.KeyCtrlR,
		gocui.ModNone,
		m.reloadmodal); err != nil {
		return err
	}

	if err := g.SetKeybinding("", '/', gocui.ModNone, toSearchView); err != nil {
		return err
	}

	return nil
}