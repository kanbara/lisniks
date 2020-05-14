package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type PartOfSpeechView struct {
	DefaultView
}

func (p *PartOfSpeechView) New(g *gocui.Gui, name string) error {

	maxX, _ := g.Size()

	if v, err := g.SetView(name, 21, 7, maxX-1, 9, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Frame = false
		err := p.Update(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PartOfSpeechView) Update(v *gocui.View) error {
	v.Clear()

	if p.State.CurrentWord() != nil {
		pos := p.Dict.PartsOfSpeech.GetByID(p.State.CurrentWord().Type)
		_, err := fmt.Fprintln(v, POSColour(pos, p.State.CurrentWord().Type))

		if err != nil {
			return err
		}
	}

	return nil
}
