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

func (s *SearchView) execSearch(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	s.state.QueuePos = -1
	s.state.CurrentSearch = ""

	var newWords lexicon.Lexicon
	search, err := v.Line(0)
	if err != nil {
		newWords = s.dict.Lexicon.Words()
		s.state.StatusText = fmt.Sprintf("")
	} else {
		newWords, err = s.dict.Lexicon.FindWords(search, s.state.SearchPattern, s.state.SearchType)
		if err != nil {
			s.state.StatusText = fmt.Sprintf("%v", err)
			if err := s.updateStatusView(g); err != nil {
				return err
			}

			return toView(g, lexView)
		}

		if err := v.SetCursor(0, 0); err != nil {
			return err
		}

		s.state.SearchQueue.Enqueue(search)
		s.state.StatusText = fmt.Sprintf("search for «%v» found %v words",
			search, len(newWords))
	}

	s.state.Words = newWords
	v.Clear()
	g.Update(func(g *gocui.Gui) error {

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
	s.state.QueuePos = -1
	s.state.CurrentSearch = ""
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

func (s *SearchView) advanceSearchType(g *gocui.Gui, _ *gocui.View) error {
	l := len(s.state.SearchTypes)
	s.state.SearchType = (s.state.SearchType + 1) % lexicon.SearchType(l)

	if v, err := g.View(searchView); err != nil {
		return nil
	} else {
		s.updateTitle(v)
	}

	return nil
}

func (s *SearchView) advanceSearchPattern(g *gocui.Gui, _ *gocui.View) error {
	l := len(s.state.SearchPatterns)
	s.state.SearchPattern = (s.state.SearchPattern + 1) % lexicon.SearchPattern(l)

	if v, err := g.View(searchView); err != nil {
		return nil
	} else {
		s.updateTitle(v)
	}

	return nil
}

func (s *SearchView) updateTitle(v *gocui.View) {
	v.Title = fmt.Sprintf("search %v %v",
		s.state.SearchTypes[s.state.SearchType],
		s.state.SearchPatterns[s.state.SearchPattern])
}

func (s *SearchView) moveQueue(g *gocui.Gui, v *gocui.View, move int) error {
	// ensure we have a queue to search through
	if s.state.SearchQueue.Len() == 0 {
		return nil
	}

	if s.state.QueuePos == -1 {
		search, err := v.Line(0)
		if err != gocui.ErrInvalidPoint && search != "" {
			s.state.CurrentSearch = search
		}
	}

	// set bounds appropriately so we don't go over or under the valid positions
	if s.state.QueuePos+move >= s.state.SearchQueue.Len() {
		return nil
	}

	if s.state.QueuePos+move < 0 {
		// pop the current search back and set the queue pos so we can save the state again

		s.state.QueuePos = -1
		v.Clear()
		// write the word at 0,0, not where the cursor was before
		if err := v.SetWritePos(0, 0); err != nil {
			return err
		}

		if err := v.SetCursor(len(s.state.CurrentSearch), 0); err != nil {
			return err
		}

		v.WriteString(s.state.CurrentSearch)

		return nil
	}

	// advance the queue position
	s.state.QueuePos = s.state.QueuePos + move

	if peek := s.state.SearchQueue.Peek(s.state.QueuePos); peek != nil {
		v.Clear()

		// write the word at 0,0, not where the cursor was before
		if err := v.SetWritePos(0, 0); err != nil {
			return err
		}

		v.WriteString(*peek)
		if err := v.SetCursor(len(*peek), 0); err != nil {
			return err
		}
	}

	return nil
}

func (s *SearchView) queueUp(g *gocui.Gui, v *gocui.View) error {
	if err := s.moveQueue(g, v, 1); err != nil {
		return err
	}

	return nil
}

func (s *SearchView) queueDown(g *gocui.Gui, v *gocui.View) error {
	if err := s.moveQueue(g, v, -1); err != nil {
		return err
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

	if err := g.SetKeybinding(lexView, 't', gocui.ModNone, s.advanceSearchType); err != nil {
		return err
	}

	if err := g.SetKeybinding(lexView, 'f', gocui.ModNone, s.advanceSearchPattern); err != nil {
		return err
	}

	if err := g.SetKeybinding(searchView, gocui.KeyArrowUp, gocui.ModNone, s.queueUp); err != nil {
		return err
	}

	if err := g.SetKeybinding(searchView, gocui.KeyArrowDown, gocui.ModNone, s.queueDown); err != nil {
		return err
	}
	return nil
}
