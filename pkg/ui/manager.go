package ui

import (
	"github.com/awesome-gocui/gocui"
	"github.com/kanbara/lisniks/pkg/dictionary"
	"github.com/kanbara/lisniks/pkg/state"
)

type Manager struct {
	dict  *dictionary.Dictionary
	state *state.State
	views map[string]ViewUpdateSetter
}

func NewManager(dict *dictionary.Dictionary, state *state.State) *Manager {
	return &Manager{dict: dict, state: state}
}

func (m *Manager) Layout(g *gocui.Gui) error {
	if m.views == nil {

		// we can share this as a singleton as these functions are all nil
		// and all the updates are thread safe anyway
		dv := DefaultView{m}

		m.views = map[string]ViewUpdateSetter{
			headerView: &HeaderView{dv},

			searchView: &SearchView{View{
				m,
				[]string{
				posView, lexView, currentWordView, localWordView,
				wordGrammarView, defnView, statusView}}},

			lexView: &LexiconView{View{
				m,
				[]string{
				posView, currentWordView, localWordView,
				wordGrammarView, defnView, statusView}}},

			currentWordView: &CurrentWordView{dv},
			localWordView:   &LocalWordView{dv},
			posView:         &PartOfSpeechView{dv},
			wordGrammarView: &WordGrammarView{dv},
			defnView:        &DefinitionView{dv},
			statusView:      &StatusView{dv},
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
