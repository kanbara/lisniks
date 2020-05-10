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

		s.updateTitle(v, search.TypeAustrianWord, search.PatternRegex)
		v.Frame = true
		v.Editable = true
		v.Editor = gocui.EditorFunc(s.updateSearchbarEditor)
	}

	return nil
}

func (s *SearchView) execSearch(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	s.state.QueuePos = -1
	s.state.CurrentSearch = ""

	var newWords lexicon.Lexicon
	word, err := v.Line(0)
	if err != nil {
		newWords = s.dict.Lexicon.Words()
		s.state.StatusText = fmt.Sprintf("")
	} else {
		newWords, err = s.dict.Lexicon.FindWords(word)
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

		s.state.SearchQueue.Enqueue(word)
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
	s.state.CurrentSearch = ""

	v.Clear()
	s.updateTitle(v, search.TypeAustrianWord, search.PatternRegex)
	s.state.StatusText = "search canceled"

	if err := s.Manager.updateStatusView(g); err != nil {
		return err
	}

	if err := v.SetCursor(0, 0); err != nil {
		return err
	}

	return toView(g, lexView)
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
			s.state.CurrentSearch = word
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
		if err := v.SetCursor(len([]rune(s.state.CurrentSearch)), 0); err != nil {
			return err
		}

		// write the word, and also pop the current search states so that we
		// end up making the correct search
		v.WriteString(s.state.CurrentSearch)
		parsed, err := search.ParseString(s.state.CurrentSearch)
		if err == nil {
			s.updateTitle(v, parsed.Type, parsed.Pattern)
		}

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
		parsed, err := search.ParseString(p)
		if err == nil {
			s.updateTitle(v, parsed.Type, parsed.Pattern)

		}

		v.WriteString(p)

		// []rune is needed here, or else we get the wrong string len with >1 byte chars!
		if err := v.SetCursor(len([]rune(p)), 0); err != nil {
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

func (s *SearchView) moveLeft(_ *gocui.Gui, v *gocui.View) error {
	if err := v.SetCursor(0, 0); err != nil {
		return err
	}

	return nil
}

func (s *SearchView) moveRight(_ *gocui.Gui, v *gocui.View) error {
	w := v.ViewBuffer()
	if w != "" {
		if err := v.SetCursor(len([]rune(w)), 0); err != nil {
			return err
		}
	}

	return nil
}

func (s *SearchView) delete(_ *gocui.Gui, v *gocui.View) error {
	v.Clear()
	if err := v.SetCursor(0, 0); err != nil {
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

	if err := g.SetKeybinding(searchView, gocui.KeyArrowUp, gocui.ModNone, s.queueUp); err != nil {
		return err
	}

	if err := g.SetKeybinding(searchView, gocui.KeyArrowDown, gocui.ModNone, s.queueDown); err != nil {
		return err
	}

	if err := g.SetKeybinding(searchView, gocui.KeyCtrlA, gocui.ModNone, s.moveLeft); err != nil {
		return err
	}

	if err := g.SetKeybinding(searchView, gocui.KeyCtrlE, gocui.ModNone, s.moveRight); err != nil {
		return err
	}

	if err := g.SetKeybinding(searchView, gocui.KeyCtrlW, gocui.ModNone, s.delete); err != nil {
		return err
	}

	return nil
}

func (s *SearchView) updateSearchbarEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	gocui.DefaultEditor.Edit(v, key, ch, mod)

	parsed, err := search.ParseString(v.Buffer())
	if err == nil { // we can't handle errors here, when the user does something bad, just ignore it
		s.updateTitle(v, parsed.Type, parsed.Pattern)
	}
}
