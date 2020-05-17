package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/kanbara/lisniks/pkg/polyglot/lexicon"
	"github.com/kanbara/lisniks/pkg/search"
)

type SearchView struct {
	View
}

func (s *SearchView) New(name string) error {

	maxX, maxY := s.g.Size()

	if v, err := s.g.SetView(name, 0, maxY-3, maxX-1, maxY-1, 0); err != nil {
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

func (s *SearchView) back() error {
	return s.ToDefaultViews([]string{POSSelectViewName, WCGSelectViewName})
}

func (s *SearchView) execSearch(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	s.State.SearchState.QueuePos = -1
	s.State.SearchState.CurrentSearch = ""

	var newWords lexicon.Lexicon
	word, err := v.Line(0)
	if err != nil {
		newWords = s.Dict.Lexicon.Words()
		s.State.StatusText = fmt.Sprintf("")
	} else {
		newWords, err = s.Dict.Lexicon.FindWords(word)
		if err != nil {
			s.State.StatusText = fmt.Sprintf("%v", err)
			if err := s.UpdateStatusView(); err != nil {
				return err
			}

			// TODO refactor this as a handler for viewPopped
			return s.back()

		}

		if err := v.SetCursor(0, 0); err != nil {
			return err
		}

		s.State.SearchState.SearchQueue.Enqueue(word)
		s.State.StatusText = fmt.Sprintf("search for «%v» found %v words",
			word, len(newWords))
	}

	s.State.Words = newWords
	s.State.SelectedWord = 0

	v.Clear()
	s.UpdateViews()

	return s.back()
}

func (s *SearchView) Update(_ *gocui.View) error { return nil }

func (s *SearchView) cancelToLexView(g *gocui.Gui, v *gocui.View) error {
	s.State.SearchState.QueuePos = -1
	s.State.SearchState.CurrentSearch = ""

	v.Clear()
	s.updateTitle(v, search.TypeAustrianWord, search.PatternRegex)
	s.State.StatusText = "search canceled"

	if err := s.UpdateStatusView(); err != nil {
		return err
	}

	return s.back()
}

func (s *SearchView) updateTitle(v *gocui.View, t search.Type, p search.Pattern) {
	title := fmt.Sprintf("search %v %v",
		s.State.SearchState.SearchTypes[t],
		s.State.SearchState.SearchPatterns[p])

	v.Title = title
}

func (s *SearchView) moveQueue(_ *gocui.Gui, v *gocui.View, move int) error {
	// ensure we have a queue to search through
	if s.State.SearchState.SearchQueue.Len() == 0 {
		return nil
	}

	if s.State.SearchState.QueuePos == -1 {
		word, err := v.Line(0)
		if err != gocui.ErrInvalidPoint && word != "" {
			s.State.SearchState.CurrentSearch = word
		}
	}

	// set bounds appropriately so we don't go over or under the valid positions
	if s.State.SearchState.QueuePos+move >= s.State.SearchState.SearchQueue.Len() {
		return nil
	}

	if s.State.SearchState.QueuePos+move < 0 {
		// pop the current search back and set the queue pos so we can save the state again

		s.State.SearchState.QueuePos = -1
		v.Clear()
		// write the word at 0,0, not where the cursor was before
		if err := v.SetWritePos(0, 0); err != nil {
			return err
		}

		// []rune is needed here, or else we get the wrong string len with >1 byte chars!
		if err := v.SetCursor(len([]rune(s.State.SearchState.CurrentSearch)), 0); err != nil {
			return err
		}

		// write the word, and also pop the current search states so that we
		// end up making the correct search
		v.WriteString(s.State.SearchState.CurrentSearch)
		parsed, err := search.ParseString(s.State.SearchState.CurrentSearch)
		if err == nil {
			s.updateTitle(v, parsed.Type, parsed.Pattern)
		}

		return nil
	}

	// advance the queue position
	s.State.SearchState.QueuePos = s.State.SearchState.QueuePos + move

	// if we have a word in the queue at this index
	if peek := s.State.SearchState.SearchQueue.Peek(s.State.SearchState.QueuePos); peek != nil {
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

func (s *SearchView) SetKeybindings() error {
	if err := s.g.SetKeybinding(SearchViewName, gocui.KeyEsc, gocui.ModNone, s.cancelToLexView); err != nil {
		return err
	}

	if err := s.g.SetKeybinding(SearchViewName, gocui.KeyEnter, gocui.ModNone, s.execSearch); err != nil {
		return err
	}

	if err := s.g.SetKeybinding(SearchViewName, gocui.KeyArrowUp, gocui.ModNone, s.queueUp); err != nil {
		return err
	}

	if err := s.g.SetKeybinding(SearchViewName, gocui.KeyArrowDown, gocui.ModNone, s.queueDown); err != nil {
		return err
	}

	if err := s.g.SetKeybinding(SearchViewName, gocui.KeyCtrlA, gocui.ModNone, s.moveLeft); err != nil {
		return err
	}

	if err := s.g.SetKeybinding(SearchViewName, gocui.KeyCtrlE, gocui.ModNone, s.moveRight); err != nil {
		return err
	}

	if err := s.g.SetKeybinding(SearchViewName, gocui.KeyCtrlW, gocui.ModNone, s.delete); err != nil {
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
