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
		v.Highlight = true

		if err := p.Update(v); err != nil {
			return err
		}
	}

	return nil
}

func (p *POSSelectView) Update(v *gocui.View) error {
	v.Clear()

	for _, n := range p.State.SearchState.POSList {
		for _, i := range p.State.SearchState.POSes {
			if int(n.ID) == i {
				if _, err := fmt.Fprintln(v, POSColourInvert(n.Name, n.ID)); err != nil {
					return err
				}
			}
		}

		// todo don't write twice
		// make pos thing map bc then it's much easier to check equality and O(1)
		// where is gender in nouns...

		if _, err := fmt.Fprintln(v, POSColour(n.Name, n.ID)); err != nil {
			return err
		}

	}

	return nil
}


func (vm *ViewManager) selectedPOS(_ *gocui.Gui, v *gocui.View) error {
	line := StripANSI(v.BufferLines()[vm.State.SearchState.SelectedPOS])
	posID := vm.Dict.PartsOfSpeech.GetByName(line)
	vm.Log.Debugf("selected line %v %v", line, posID)

	err := v.SetLine(vm.State.SearchState.SelectedPOS, POSColourInvert(line, posID))
	if err != nil {
		return err
	}

	// selected POS should be map prolly for easier add remove than slice
	vm.State.SearchState.POSes = append(vm.State.SearchState.POSes, int(posID))

	// TODO
	// add to state
	// remove to state
	// invert selection?
	// todo in listview, if bold show as unbold on highlight ? maybe 256 instead.. add tix

	return nil
}