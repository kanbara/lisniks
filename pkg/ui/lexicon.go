package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

func (m *Manager) NewLexiconView(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView(lexView, 0, 3, 20, maxY-1, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = "lexicon"
		v.Frame = true
		v.Highlight = true
		v.SelFgColor = gocui.ColorGreen
		v.FgColor = gocui.ColorWhite

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

		err = v.SetHighlight(m.state.SelectedWord, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) updateWord(g *gocui.Gui, v *gocui.View, updown int) error {
	// doesn't scroll or wrap

	x, y := v.Cursor()
	//if y+updown < 0 {
	//	return nil
	//}
	//
	//_, sizeY := v.Size()
	//if y+updown >= sizeY {
	//	//v.SetOrigin(x)
	//	return nil
	//}

	err :=v.SetCursor(x, y+updown)
	if err != nil {
		return err
	}

	// turn off highlight for previous words
	err = v.SetHighlight(m.state.SelectedWord, false)
	if err != nil {
		return err
	}

	m.state.SelectedWord = m.state.SelectedWord + updown

	// highlight current word
	err = v.SetHighlight(m.state.SelectedWord, true)
	if err != nil {
		return err
	}

	g.Update(func(g *gocui.Gui) error {
		// TODO loop thru the views
		// also make View have an Update function so we can use it with interface
		// and then call update... oy vey
		v, err := g.View(posView)
		if err != nil {
			return err
		}

		err = m.UpdatePartOfSpeech(v)
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