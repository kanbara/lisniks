package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type CurrentWordView struct {
	DefaultView
}

func (c *CurrentWordView) New(name string) error {
	if v, err := c.g.SetView(name, 0, 3, 20, 5, 0); err != nil {
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

	if c.State.CurrentWord() != nil {
		_, err := fmt.Fprintln(v, c.State.CurrentWord().Austrian)
		if err != nil {
			return err
		}
	}

	return nil
}
