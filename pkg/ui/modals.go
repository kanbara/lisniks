package ui

import (
	"github.com/awesome-gocui/gocui"
	"github.com/kanbara/lisniks/pkg/dictionary"
	"github.com/kanbara/lisniks/pkg/state"
)

// used to send the `error` which quits the program
func (vm *ViewManager) QuitModal(g *gocui.Gui, _ *gocui.View) error {
	AddModalView(g, vm, "quit the program?", func(_ *gocui.Gui, _ *gocui.View) error {
		return gocui.ErrQuit
	}, ModalQuit)

	return nil
}

func (vm *ViewManager) ReloadModal(g *gocui.Gui, _ *gocui.View) error {
	AddModalView(g, vm, "reload dictionary?", func(_ *gocui.Gui, v *gocui.View) error {

		// replace State and Dict
		dict := dictionary.NewDictFromFile(vm.Dict.Filename(), vm.Log)
		s := state.NewState(vm.State.Version, dict)

		vm.Dict = dict
		vm.State = s

		// update all Views
		for name := range vm.Views {
			// this is THAT view by name (e.g. headerView itself) not this view
			if viewsV, err := g.View(name); err != nil {
				return err
			} else {
				if err := vm.Views[name].Update(viewsV); err != nil {
					return err
				}
			}
		}

		return nil
	}, ModalReload)

	return nil
}
