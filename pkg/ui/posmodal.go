package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type POSSelectView struct {
	ListView
}

func (p *POSSelectView) New(name string) error {
	x, y := p.g.Size()
	if v, err := p.g.SetView(name, 5, 5, (x/2)-5, y-5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = name
		v.Frame = true
		v.FgColor = gocui.ColorRed

		err := p.Update(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *POSSelectView) Update(v *gocui.View) error {
	v.Clear()

	m := p.Dict.PartsOfSpeech.GetNameToIDs()

	for n, t := range m {
		_, err := fmt.Fprintln(v, POSColour(n, t))

		if err != nil {
			return err
		}
	}

	return nil
}
