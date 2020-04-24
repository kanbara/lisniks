package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func (m *Manager) NewLexiconView(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView(lexView, 0, 3, 20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "lexicon"
		v.Frame = true
		v.Highlight = true
		v.SelBgColor = gocui.ColorMagenta
		v.SelFgColor = gocui.ColorBlack

		for _, w  := range m.state.Words {
			_, err := fmt.Fprintln(v, w.Con)
			if err != nil {
				return err
			}
		}

		err := v.SetOrigin(0, 0)
		if err != nil {
			return err
		}

		err = v.SetCursor(0, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) updateWord(g *gocui.Gui, v *gocui.View, updown int) error {
	// doesn't scroll or wrap
	x, y := v.Cursor()
	err :=v.SetCursor(x, y+updown)
	if err != nil {
		return err
	}

	m.state.SelectedWord = m.state.SelectedWord + updown
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(posView)
		if err != nil {
			return err
		}

		err = m.UpdatePartOfSpeech(v)
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

	return nil
}

func (m *Manager) nextWord(g *gocui.Gui, v *gocui.View) error {
	err := m.updateWord(g, v, 1)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) nextWordJump(g *gocui.Gui, v *gocui.View) error {
	err := m.updateWord(g, v, 10)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) prevWord(g *gocui.Gui, v *gocui.View) error {
	err := m.updateWord(g, v, -1)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) prevWordJump(g *gocui.Gui, v *gocui.View) error {
	err := m.updateWord(g, v, -10)
	if err != nil {
		return err
	}

	return nil
}