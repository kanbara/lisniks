package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type DefinitionView NoBindingsView

func (d *DefinitionView) New(g *gocui.Gui, name string) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView(name, 21, 10, maxX-1, maxY-4, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = name
		v.Wrap = true
		err := d.Update(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DefinitionView) Update(v *gocui.View) error {
	v.Clear()

	if d.state.CurrentWord() != nil {
		_, err := fmt.Fprintln(v, d.state.CurrentWord().Def)
		if err != nil {
			return err
		}
	}

	return nil
}
