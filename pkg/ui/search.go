package ui

import (
	"github.com/awesome-gocui/gocui"
	"github.com/kanbara/lisniks/pkg/lexicon"
)

func (m *Manager) NewSearchView(g* gocui.Gui) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView(searchView, 0, maxY-3, maxX-1, maxY-1, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = searchView
		v.Frame = true
		v.Editable = true
	}

	return nil
}

// todo add flag and field for fuzzy
// todo add regex
// todo add statusbar showing # matches found, time, flags

func (m *Manager) execSearch(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	err := v.SetCursor(0,0)
	if err != nil {
		return err
	}

	var newWords lexicon.Lexicon
	search, err := v.Line(0)
	if err != nil {
		newWords = m.dict.Lexicon.Words()
	} else {
		newWords = m.dict.Lexicon.FindByConWordFuzzy(search)
	}

	g.Update(func(g *gocui.Gui) error {
		// todo fixme
		m.state.Words = newWords
		v.Clear()

		v, err := g.View(lexView)
		if err != nil {
			return err
		}

		err = m.UpdateLexicon(v)
		if err != nil {
			return err
		}

		v, err = g.View(posView)
		if err != nil {
			return err
		}

		err = m.UpdatePartOfSpeech(v)
		if err != nil {
			return err
		}

		v, err = g.View(currentWordView)
		if err != nil {
			return err
		}

		err = m.updateCurrentWordView(v)
		if err != nil {
			return err
		}

		v, err = g.View(localWordView)
		if err != nil {
			return err
		}

		err = m.UpdateLocalWordView(v)
		if err != nil {
			return err
		}

		v, err = g.View(wordGrammarView)
		if err != nil {
			return err
		}

		err = m.UpdateWordGrammarView(v)
		if err != nil {
			return err
		}

		v, err = g.View(defnView)
		if err != nil {
			return err
		}

		err = m.UpdateDefinition(v)
		if err != nil {
			return err
		}

		return nil
	})

	return toView(g, lexView)
}