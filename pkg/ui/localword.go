package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type LocalWordView DefaultView

func (l *LocalWordView) New(g *gocui.Gui, name string) error {

	maxX, _ := g.Size()

	if v, err := g.SetView(name, 21, 3, maxX-1, 5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = name
		v.FgColor = gocui.ColorYellow
		err := l.Update(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LocalWordView) Update(v *gocui.View) error {
	v.Clear()

	if l.state.CurrentWord() != nil {
		_, err := fmt.Fprintln(v, l.state.CurrentWord().Local)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LocalWordView) SetKeybindings(_ *gocui.Gui) error { return nil }