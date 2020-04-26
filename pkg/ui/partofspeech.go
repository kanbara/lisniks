package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

func (m *Manager) NewPartOfSpeechView(g* gocui.Gui) error {

	maxX, _ := g.Size()

	if v, err := g.SetView(posView, 21, 11, maxX-1, 13, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = posView
		v.FgColor = colour(int(m.state.CurrentWord().Type))
		err := m.UpdatePartOfSpeech(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) UpdatePartOfSpeech(v *gocui.View) error {
	v.Clear()
	v.FgColor = colour(int(m.state.CurrentWord().Type))

	_, err := fmt.Fprintln(v, m.dict.PartsOfSpeech.Get(m.state.CurrentWord().Type))
	if err != nil {
		return err
	}

	return nil
}