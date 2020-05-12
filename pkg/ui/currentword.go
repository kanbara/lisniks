package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type CurrentWordView struct {
	DefaultView
}

func (c *CurrentWordView) New(g *gocui.Gui, name string) error {
	if v, err := g.SetView(name, 0, 3, 20, 5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Frame = true
		v.FgColor = gocui.ColorGreen

		err := c.Update(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CurrentWordView) Update(v *gocui.View) error {
	v.Clear()

	if c.state.CurrentWord() != nil {
		_, err := fmt.Fprintln(v, c.state.CurrentWord().Austrian)
		if err != nil {
			return err
		}
	}

	return nil
}
