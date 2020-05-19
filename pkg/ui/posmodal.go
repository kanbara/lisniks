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
		v.Highlight = true

		// todo we need this rn because we create the view again but don't reset the sel
		p.State.SearchState.SelectedPOS = 0

		if err := p.Update(v); err != nil {
			return err
		}
	}

	return nil
}

func (p *POSSelectView) Update(v *gocui.View) error {
	v.Clear()

	for _, n := range p.State.SearchState.POSList {
		if isSet, ok := p.State.SearchState.POSes[int(n.ID)]; ok && isSet {
			if _, err := fmt.Fprintln(v, POSColourInvert(n.Name, n.ID)); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintln(v, POSColour(n.Name, n.ID)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *ViewManager) selectedPOS(_ *gocui.Gui, v *gocui.View) error {
	line := StripANSI(v.BufferLines()[vm.State.SearchState.SelectedPOS])
	posID := vm.Dict.PartsOfSpeech.GetByName(line)
	vm.Log.Debugf("selected line %v %v", line, posID)

	if isSet, ok := vm.State.SearchState.POSes[int(posID)]; ok {
		vm.State.SearchState.POSes[int(posID)] = !isSet
		if isSet {
			if err := v.SetLine(vm.State.SearchState.SelectedPOS, POSColour(line, posID)); err != nil {
				return err
			}
		} else {
			if err := v.SetLine(vm.State.SearchState.SelectedPOS, POSColourInvert(line, posID)); err != nil {
				return err
			}
		}
	}

	// todo in listview, if bold show as unbold on highlight ? maybe 256 instead.. add tix

	return nil
}
