package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type HeaderView DefaultView

func (h *HeaderView) New(g *gocui.Gui, name string) error {
	stats := h.dict.Stats()
	langAndVersion := h.dict.LangAndVersion()
	stringlen := len(stats)

	// TODO move positions where possible into views.go maybe
	if v, err := g.SetView(name, 0, 0, stringlen+1, 5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		_, err := fmt.Fprintln(v, langAndVersion)
		_, err = fmt.Fprintln(v, stats)

		if err != nil {
			return err
		}

		v.Frame = false
	}

	return nil
}

func (h *HeaderView) Update(_ *gocui.View) error { return nil }