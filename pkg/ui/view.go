package ui

import (
	"github.com/awesome-gocui/gocui"
)

type ViewUpdateSetter interface {
	New(g *gocui.Gui, name string) error
	Update(v *gocui.View) error
	SetKeybindings(g *gocui.Gui) error
}

type View struct {
	*ViewManager
	ViewsToUpdate []string
}

func (vw *View) UpdateViews(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		for _, viewName := range vw.ViewsToUpdate {
			if v, err := g.View(viewName); err != nil {
				return err
			} else {
				if err := vw.Views[viewName].Update(v); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

type ListView struct {
	View
}

type DefaultView struct {
	*ViewManager
}

func (d *DefaultView) Update(_ *gocui.View) error        { return nil }
func (d *DefaultView) SetKeybindings(_ *gocui.Gui) error { return nil }