package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"strings"
)

type modalType string

const (
	modalQuit   modalType = "quit"
	modalReload           = "reload"
)

func modalName(mt modalType) string {
	return string(mt) + "modal"
}

type ModalView struct {
	DefaultView
	text             string
	action           func(g *gocui.Gui, v *gocui.View) error
	mt               modalType
	currentView      *gocui.View // for going to upon cancel
	frameColour      gocui.Attribute
	frameTitleColour gocui.Attribute
}

func (m *Manager) AddModalView(g *gocui.Gui,
	text string,
	action func(g *gocui.Gui, v *gocui.View) error,
	mt modalType) { // TODO add keybinding here to blacklist and whitelist after
	g.Update(func(g *gocui.Gui) error {
		mv := ModalView{DefaultView: DefaultView{m, m.log},
			text:   text,
			action: action,
			mt:     mt}

		if err := mv.New(g, modalName(mt)); err != nil {
			return err
		}

		mv.currentView = g.CurrentView()

		if _, err := g.SetCurrentView(modalName(mt)); err != nil {
			return err
		}

		// disable all global keybindings to not f**k things up
		for _, k := range GlobalKeybindingKeys() {
			if err := g.BlacklistKeybinding(k); err != nil {
				return err
			}
		}

		return nil
	})
}

func (m *ModalView) New(g *gocui.Gui, name string) error {
	x, y := g.Size()
	if v, err := g.SetView(name,
		x/2-20, y/2-5,
		x/2+20, y/2+5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Frame = true

		m.frameColour = g.SelFrameColor
		m.frameTitleColour = g.SelFgColor

		g.SelFrameColor = gocui.ColorBlue
		g.SelFgColor = gocui.ColorBlue

		// proper sizes, ugly but works
		x, y := v.Size()
		midX := (x - len(m.text)) / 2
		midY := (y - 1) / 2
		ystr := strings.Repeat("\n", midY)
		xstr := strings.Repeat(" ", midX)

		if _, err := fmt.Fprintln(v, ystr+xstr+m.text); err != nil {
			return err
		}

		opText := fmt.Sprintf("%v or %v",
			applyColour(Green, "[ENTER] Yes", ansibold),
			applyColour(Red, "[ESC] Cancel", ansibold))

		l := (x - len(stripANSI(opText)))/2
		if l <= 0 {
			l = 0
		}

		if _, err := fmt.Fprintln(v,
			strings.Repeat(" ", l)+opText); err != nil {
			return err
		}

		if err := m.SetKeybindings(g); err != nil {
			return err
		}
	}

	return nil
}

func (m *ModalView) cleanup(g *gocui.Gui) error {
	if err := g.DeleteView(modalName(m.mt)); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(m.currentView.Name()); err != nil {
		return err
	}

	g.SelFrameColor = m.frameColour
	g.SelFgColor = m.frameTitleColour

	for _, k := range GlobalKeybindingKeys() {
		if err := g.WhitelistKeybinding(k); err != nil {
			return err
		}
	}

	return nil
}

func (m *ModalView) SetKeybindings(g *gocui.Gui) error {

	//success func
	if err := g.SetKeybinding(modalName(m.mt),
		gocui.KeyEnter,
		gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			// execute the action for the keybinding and then cleanup
			if err := m.action(g, v); err != nil {
				return err
			}

			return m.cleanup(g)
		}); err != nil {
		return err
	}

	// cancel func
	if err := g.SetKeybinding(modalName(m.mt),
		gocui.KeyEsc,
		gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			// basically do nothing here
			return m.cleanup(g)
		}); err != nil {
		return err
	}

	return nil
}
