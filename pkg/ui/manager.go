package ui

import (
	"github.com/awesome-gocui/gocui"
	"github.com/kanbara/lisniks/pkg/dictionary"
	"github.com/kanbara/lisniks/pkg/state"
)

type Manager struct {
	dict  *dictionary.Dictionary
	state *state.State
	views map[string]View
}

func NewManager(dict *dictionary.Dictionary, state *state.State) *Manager {
	return &Manager{dict: dict, state: state}
}

func (m *Manager) Layout(g *gocui.Gui) error {
	if m.views == nil {

		nbu := &NilBindingsAndUpdates{}
		nb := &NilBindings{}

		m.views = map[string]View{
			headerView: &HeaderView{m, nbu},
			searchView: &SearchView{m, []string{
				posView, lexView, currentWordView, localWordView, wordGrammarView, defnView}},
			lexView: &LexiconView{m, []string{
				posView, currentWordView, localWordView, wordGrammarView, defnView}},
			currentWordView: &CurrentWordView{m, nb},
			localWordView:   &LocalWordView{m, nb},
			posView:         &PartOfSpeechView{m, nb},
			wordGrammarView: &WordGrammarView{m, nb},
			defnView:        &DefinitionView{m, nb},
		}

		for name, view := range m.views {
			if err := view.New(g, name); err != nil {
				return err
			}

			if err := view.SetKeybindings(g); err != nil {
				return err
			}
		}
	}

	return nil
}
