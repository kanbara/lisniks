package ui

import (
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

		v.Title = name
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
	} else {
		newWords = s.dict.Lexicon.FindByConWordFuzzy(search)
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

func cancelToLexView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	v.Clear()
	if err := v.SetCursor(0, 0); err != nil {
		return err
	}

	return toView(g, lexView)
}

func (s *SearchView) SetKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(searchView, gocui.KeyEsc, gocui.ModNone, cancelToLexView); err != nil {
		return err
	}

	if err := g.SetKeybinding(searchView, gocui.KeyEnter, gocui.ModNone, s.execSearch); err != nil {
		return err
	}

	return nil
}
