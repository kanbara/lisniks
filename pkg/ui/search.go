package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/kanbara/lisniks/pkg/lexicon"
	"github.com/kanbara/lisniks/pkg/search"
)

type SearchView DefaultView

func (s *SearchView) New(g *gocui.Gui, name string) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView(name, 0, maxY-3, maxX-1, maxY-1, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		s.updateTitle(v, s.state.SearchType, s.state.SearchPattern)
		v.Frame = true
		v.Editable = true
	}

	return nil
}

func (s *SearchView) execSearch(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	s.state.QueuePos = -1
	s.state.CurrentSearch = search.Data{}

	var newWords lexicon.Lexicon
	word, err := v.Line(0)
	if err != nil {
		newWords = s.dict.Lexicon.Words()
		s.state.StatusText = fmt.Sprintf("")
	} else {
		newWords, err = s.dict.Lexicon.FindWords(word, s.state.SearchPattern, s.state.SearchType)
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

		s.state.SearchQueue.Enqueue(search.Data{
			Type: s.state.SearchType,
			Pattern: s.state.SearchPattern,
			String: word,
		})
		s.state.StatusText = fmt.Sprintf("search for «%v» found %v words",
			word, len(newWords))
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
	s.state.CurrentSearch = search.Data{}

	v.Clear()
	s.updateTitle(v, s.state.SearchType, s.state.SearchPattern)
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
	s.state.SearchType = (s.state.SearchType + 1) % search.Type(l)

	if v, err := g.View(searchView); err != nil {
		return nil
	} else {
		s.updateTitle(v, s.state.SearchType, s.state.SearchPattern)
	}

	return nil
}

func (s *SearchView) advanceSearchPattern(g *gocui.Gui, _ *gocui.View) error {
	l := len(s.state.SearchPatterns)
	s.state.SearchPattern = (s.state.SearchPattern + 1) % search.Pattern(l)

	if v, err := g.View(searchView); err != nil {
		return nil
	} else {
		s.updateTitle(v, s.state.SearchType, s.state.SearchPattern)
	}

	return nil
}

func (s *SearchView) updateTitle(v *gocui.View, t search.Type, p search.Pattern) {
	v.Title = fmt.Sprintf("search %v %v",
		s.state.SearchTypes[t],
		s.state.SearchPatterns[p])
}

func (s *SearchView) moveQueue(g *gocui.Gui, v *gocui.View, move int) error {
	// ensure we have a queue to search through
	if s.state.SearchQueue.Len() == 0 {
		return nil
	}

	if s.state.QueuePos == -1 {
		word, err := v.Line(0)
		if err != gocui.ErrInvalidPoint && word != "" {
			s.state.CurrentSearch = search.Data{
				Type: s.state.SearchType,
				Pattern: s.state.SearchPattern,
				String: word,
			}
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

		// []rune is needed here, or else we get the wrong string len with >1 byte chars!
		if err := v.SetCursor(len([]rune(s.state.CurrentSearch.String)), 0); err != nil {
			return err
		}

		// write the word, and also pop the current search states so that we
		// end up making the correct search
		v.WriteString(s.state.CurrentSearch.String)
		s.updateTitle(v, s.state.CurrentSearch.Type, s.state.CurrentSearch.Pattern)
		s.state.SearchType = s.state.CurrentSearch.Type
		s.state.SearchPattern = s.state.CurrentSearch.Pattern
		return nil
	}

	// advance the queue position
	s.state.QueuePos = s.state.QueuePos + move

	// if we have a word in the queue at this index
	if peek := s.state.SearchQueue.Peek(s.state.QueuePos); peek != nil {
		v.Clear()

		p := *peek
		// write the word at 0,0, not where the cursor was before
		if err := v.SetWritePos(0, 0); err != nil {
			return err
		}

		// write the word and pop the search states here too!
		s.state.SearchType = p.Type
		s.state.SearchPattern = p.Pattern
		s.updateTitle(v, p.Type, p.Pattern)
		v.WriteString(p.String)

		// []rune is needed here, or else we get the wrong string len with >1 byte chars!
		if err := v.SetCursor(len([]rune(p.String)), 0); err != nil {
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
