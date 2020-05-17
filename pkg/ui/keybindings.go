package ui

import (
	"github.com/awesome-gocui/gocui"
)

func GlobalKeybindingKeys() []gocui.Key {
	// XXX that's dumb
	// the black/whitelist code only checks the KEY but for runes
	// the KEY is always 0, had to scour the source for this. whatever
	return []gocui.Key{gocui.KeyCtrlR, gocui.KeyCtrlC, 0, gocui.KeyTab}
}