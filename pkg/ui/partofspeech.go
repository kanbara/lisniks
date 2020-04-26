package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type PartOfSpeechView NoBindingsView

func (p *PartOfSpeechView) New(g *gocui.Gui, name string) error {

	maxX, _ := g.Size()

	if v, err := g.SetView(name, 21, 7, maxX-1, 9, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Frame = false
		v.FgColor = colour(int(p.state.CurrentWord().Type))
		err := p.Update(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PartOfSpeechView) Update(v *gocui.View) error {
	v.Clear()

	if p.state.CurrentWord() != nil {
		v.FgColor = colour(int(p.state.CurrentWord().Type))

		_, err := fmt.Fprintln(v, p.dict.PartsOfSpeech.Get(p.state.CurrentWord().Type))
		if err != nil {
			return err
		}
	}

	return nil
}
