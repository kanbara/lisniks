package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type CurrentWordView DefaultView

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
		v.FgColor = gocui.ColorGreen
		_, err := fmt.Fprintln(v, c.state.CurrentWord().Con)
		if err != nil {
			return err
		}
	} else {
		v.FgColor = gocui.ColorRed
		_, err := fmt.Fprintln(v, "NO WORDS FOUND")
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CurrentWordView) SetKeyBindings(_ *gocui.Gui) error { return nil }