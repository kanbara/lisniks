package ui

import (
	"github.com/awesome-gocui/gocui"
	"github.com/sirupsen/logrus"
)

const (
	headerView      = "lang"
	lexView         = "lexicon"
	posView         = "part of speech"
	wordGrammarView = "word classes"
	localWordView   = "local word"
	defnView        = "definition"
	currentWordView = "current word"
	searchView      = "search"
	statusView      = "status"
	debugView       = "debug"
)

type ViewUpdateSetter interface {
	New(g *gocui.Gui, name string) error
	Update(v *gocui.View) error
	SetKeybindings(g *gocui.Gui) error
}

type View struct {
	*Manager
	viewsToUpdate []string
}

type DefaultView struct {
	*Manager
	log *logrus.Logger
}

func (d *DefaultView) Update(_ *gocui.View) error { return nil }
func (d *DefaultView) SetKeybindings(_ *gocui.Gui) error { return nil }


func toSearchView(g *gocui.Gui, _ *gocui.View) error {
	g.Cursor = true
	return toView(g, searchView)
}

func toView(g *gocui.Gui, view string) error {
	_, err := g.SetCurrentView(view)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) updateStatusView(g *gocui.Gui) error {
	g.Update(func(g *gocui.Gui) error {
		if v, err := g.View(statusView); err != nil {
			return err
		} else {
			if err := m.views[statusView].Update(v); err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}
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
//	_, err := g.SetCurrentView(lexView)
//	return err
//}
