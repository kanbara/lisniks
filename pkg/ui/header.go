package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type HeaderView struct {
	DefaultView
}

func (h *HeaderView) New(g *gocui.Gui, name string) error {
	stringlen := len(h.dict.Stats())

	// TODO move positions where possible into views.go maybe
	if v, err := g.SetView(name, 0, 0, stringlen+1, 5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Frame = false
		if _, err := g.SetViewOnBottom(name); err != nil {
			return err
		}

		if err := h.Update(v); err != nil {
			return err
		}
	}


	return nil
}

func (h *HeaderView) Update(v *gocui.View) error {
	v.Clear()


	if _, err := fmt.Fprintln(v, fmt.Sprintf("%v 💛 lisniks %v",
		h.dict.LangAndVersion(), h.state.Version)); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(v, h.dict.Stats()); err != nil {
		return err
	}


	return nil
}