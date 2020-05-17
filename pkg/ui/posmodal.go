package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"sort"
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
		v.Highlight = true

		if err := p.Update(v); err != nil {
			return err
		}
	}

	return nil
}

func (p *POSSelectView) Update(v *gocui.View) error {
	v.Clear()

	names := p.Dict.PartsOfSpeech.GetNameToIDs()
	sort.Sort(names)

	for _, n := range names {
		_, err := fmt.Fprintln(v, POSColour(n.Name, n.ID))

		if err != nil {
			return err
		}
	}

	//p.UpdateViews()

	return nil
}
