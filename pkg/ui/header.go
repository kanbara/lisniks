package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type HeaderView NoBindingsOrUpdatesView

func (h *HeaderView) New(g *gocui.Gui, name string) error {
	stats := h.dict.Stats()
	langAndVersion := h.dict.LangAndVersion()
	stringlen := len(stats)

	// TODO move positions where possible into views.go maybe
	if v, err := g.SetView(name, 0, 0, stringlen+1, 5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		if _, err := g.SetViewOnBottom(name); err != nil {
			return err
		}

		_, err := fmt.Fprintln(v, fmt.Sprintf("%v 💛 lisniks %v", langAndVersion, h.state.Version))
		_, err = fmt.Fprintln(v, stats)

		if err != nil {
			return err
		}

		v.Frame = false
	}

	return nil
}
