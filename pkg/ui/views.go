package ui

import "github.com/awesome-gocui/gocui"

const (
	headerView = "lang"
	lexView = "lexicon"
	posView = "part of speech"
	wordGrammarView = "word classes"
	localWordView = "local word"
	defnView = "definition"
	currentWordView = "current word"
	searchView = "search"
)

type View interface {
	SetKeybindings(g *gocui.Gui) error
	New(g *gocui.Gui, name string) error
	Update(v *gocui.View) error
}

type DefaultView struct {
	*Manager
	viewsToUpdate []string
}

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