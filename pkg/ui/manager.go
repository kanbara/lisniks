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

// used to send the `error` which quits the program
func (m *Manager) quitmodal(g *gocui.Gui, _ *gocui.View) error {
	m.AddModalView(g, "quit the program?", func(_ *gocui.Gui, _ *gocui.View) error {
		return gocui.ErrQuit
	}, modalQuit)

	return nil
}

func (m *Manager) reloadmodal(g *gocui.Gui, _ *gocui.View) error {
	m.AddModalView(g, "reload dictionary?", func(_ *gocui.Gui, v *gocui.View) error {

		// replace state and dict
		dict := dictionary.NewDictFromFile(m.dict.Filename())
		s := state.NewState(m.state.Version, dict)

		m.dict = dict
		m.state = s

		// update all views
		for name := range m.views {
			// this is THAT view by name (e.g. headerView itself) not this view
			if viewsV, err := g.View(name); err != nil {
				return err
			} else {
				if err := m.views[name].Update(viewsV); err != nil {
					return err
				}
			}
		}

		return nil
	}, modalReload)

	return nil
}
