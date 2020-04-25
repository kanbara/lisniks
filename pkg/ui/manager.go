package ui

import (
	"github.com/awesome-gocui/gocui"
	"github.com/kanbara/lisniks/pkg/dictionary"
	 "github.com/kanbara/lisniks/pkg/state"
)

type Manager struct {
	dict *dictionary.Dictionary
	state *state.State
}

func NewManager(dict *dictionary.Dictionary, state *state.State) *Manager {
	return &Manager{dict:dict, state:state}
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

	err = m.NewLocalWordView(g)
	if err != nil {
		return nil
	}

	err = m.NewPartOfSpeechView(g)
	if err != nil {
		return nil
	}

	err = m.NewWordGrammarView(g)
	if err != nil {
		return nil
	}

	err = m.NewDefinitionView(g)
	if err != nil {
		return nil
	}

	// TODO this changes to a function when we can switch views
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.SelFrameColor = gocui.ColorGreen
	_, err = g.SetCurrentView(lexView)
	if err != nil {
		return nil
	}


	return nil
}