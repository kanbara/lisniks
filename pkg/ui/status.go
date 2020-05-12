package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type StatusView struct {
	DefaultView
}

func (s *StatusView) New(g *gocui.Gui, name string) error {

	maxX, maxY := g.Size()
	if v, err := g.SetView(name, 0, maxY-5, maxX, maxY-3, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Frame = false
		v.FgColor = gocui.ColorCyan | gocui.AttrBold

		if err := s.Update(v); err != nil {
			return err
		}
	}

	return nil
}

func (s *StatusView) Update(v *gocui.View) error {
	v.Clear()

	if _, err := fmt.Fprintln(v, s.state.StatusText); err != nil {
		return err
	}

	return nil
}
