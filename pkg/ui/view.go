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

type ListView struct {
	View
}

type DefaultView struct {
	*ViewManager
}

func (d *DefaultView) Update(_ *gocui.View) error { return nil }
func (d *DefaultView) SetKeybindings(_ *gocui.Gui) error { return nil }

//func (m *Manager) NextView(g *gocui.Gui, v *gocui.View) error {
//	a.viewIndex = (a.viewIndex + 1) % len(VIEWS)
//	return a.setView(g)
//}
//
//func (m *Manager) PrevView(g *gocui.Gui, v *gocui.View) error {
//	a.viewIndex = (a.viewIndex - 1 + len(VIEWS)) % len(VIEWS)
//	return a.setView(g)
//}

//func (m *Manager) setView(g *gocui.Gui) error {
//	_, err := g.SetCurrentView(LexView)
//	return err
//}
