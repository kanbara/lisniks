package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func (m *Manager) NewPartOfSpeechView(g* gocui.Gui) error {

	maxX, _ := g.Size()

	if v, err := g.SetView(posView, 21, 6, maxX-1, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = posView
		err := m.UpdatePartOfSpeech(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) UpdatePartOfSpeech(v *gocui.View) error {
	v.Clear()
	_, err := fmt.Fprintln(v, m.dict.PartsOfSpeech.Get(m.state.CurrentWord().Type))
	if err != nil {
		return err
	}

	return nil
}