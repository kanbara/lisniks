package ui

import (
	"github.com/jroimartin/gocui"
	"github.com/kanbara/lisniks/pkg/dictionary"
)

type Manager struct {
	dict *dictionary.Dictionary
}

func NewManager(dict *dictionary.Dictionary) *Manager {
	return &Manager{dict:dict}
}

func (m *Manager) Layout(g *gocui.Gui) error {

	// TODO fix the instantiation and object creation here
	// i would like to have some sort of map for the positions with percentages
	// as well as various function generators or sth so it's not so messy
	err := m.NewHeaderView(g)
	if err != nil {
		return nil
	}

	err = m.NewLexiconView(g)
	if err != nil {
		return nil
	}

	return nil
}