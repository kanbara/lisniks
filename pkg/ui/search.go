package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/kanbara/lisniks/pkg/lexicon"
)

type SearchView DefaultView

func (s *SearchView) New(g *gocui.Gui, name string) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView(name, 0, maxY-3, maxX-1, maxY-1, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		s.updateTitle(v)
		v.Frame = true
		v.Editable = true
	}

	return nil
}

// todo add flag and field for fuzzy
// todo add regex
// todo add statusbar showing # matches found, time, flags
func (s *SearchView) execSearch(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	if err := v.SetCursor(0, 0); err != nil {
		return err
	}

	var newWords lexicon.Lexicon
	search, err := v.Line(0)
	if err != nil {
		newWords = s.dict.Lexicon.Words()
		s.state.StatusText = fmt.Sprintf("")
	} else {
		switch s.state.SearchType {
		case lexicon.SearchTypeLocalWord:
			newWords = s.dict.Lexicon.FindLocalWords(search, s.state.SearchFuzzy)
		case lexicon.SearchTypeConWord:
			newWords = s.dict.Lexicon.FindConWords(search, s.state.SearchFuzzy)
		}
		s.state.StatusText = fmt.Sprintf("search for «%v» found %v words",
			search, len(newWords))
	}

	v.Clear()
	g.Update(func(g *gocui.Gui) error {
		s.state.Words = newWords

		for _, viewName := range s.viewsToUpdate {
			if v, err := g.View(viewName); err != nil {
				return err
			} else {
				if err := s.views[viewName].Update(v); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return toView(g, lexView)
}

func (s *SearchView) Update(_ *gocui.View) error { return nil }

func (s *SearchView) cancelToLexView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	v.Clear()
	s.state.StatusText = "search canceled"

	if err := s.Manager.updateStatusView(g); err != nil {
		return err
	}

	if err := v.SetCursor(0, 0); err != nil {
		return err
	}

	return toView(g, lexView)
}

func (s *SearchView) advanceSearchMode(g *gocui.Gui, _ *gocui.View) error {
	l := len(s.state.SearchTypes)
	s.state.SearchType = (s.state.SearchType + 1) % l

	if v, err := g.View(searchView); err != nil {
		return nil
	} else {
		s.updateTitle(v)
	}

	return nil
}

func (s *SearchView) updateTitle(v *gocui.View) {
	v.Title = fmt.Sprintf("search %v", s.state.SearchTypes[s.state.SearchType])
	if s.state.SearchFuzzy {
		v.Title += " fuzzy"
	}
}

func (s *SearchView) toggleFuzzy(g *gocui.Gui, _ *gocui.View) error {
	s.state.SearchFuzzy = !s.state.SearchFuzzy

	if v, err := g.View(searchView); err != nil {
		return nil
	} else {
		s.updateTitle(v)
	}

	return nil
}

func (s *SearchView) SetKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(searchView, gocui.KeyEsc, gocui.ModNone, s.cancelToLexView); err != nil {
		return err
	}

	if err := g.SetKeybinding(searchView, gocui.KeyEnter, gocui.ModNone, s.execSearch); err != nil {
		return err
	}

	if err := g.SetKeybinding(lexView, 't', gocui.ModNone, s.advanceSearchMode); err != nil {
		return err
	}

	if err := g.SetKeybinding(lexView, 'f', gocui.ModNone, s.toggleFuzzy); err != nil {
		return err
	}

	return nil
}
