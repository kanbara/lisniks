package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type LexiconView struct {
	ListView
}

func (l *LexiconView) New(g *gocui.Gui, name string) error {
	_, maxY := g.Size()
	if v, err := g.SetView(name, 0, 6, 20, maxY-6, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		_, err = g.SetCurrentView(name)
		if err != nil {
		}

		// TODO i wanted to set the title to the current word
		// but unfortunately it didn't work because it printed
		// sth. like -jeř-a because of the two byte width unicode
		//
		// i saw an issue on the thing about unicode, maybe there's a fix
		// that can be done.
		v.Title = fmt.Sprintf("%v %v/%v", name, len(l.State.Words), l.Dict.Lexicon.Len())
		v.Frame = true
		v.Highlight = true
		v.FgColor = gocui.ColorWhite

		for _, w := range l.State.Words {
			_, err := fmt.Fprintln(v, w.Austrian)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *LexiconView) Update(v *gocui.View) error {
	v.Clear()
	v.Title = fmt.Sprintf("%v %v/%v", LexViewName, len(l.State.Words), l.Dict.Lexicon.Len())

	if len(l.State.Words) > 0 {
		for _, w := range l.State.Words {
			_, err := fmt.Fprintln(v, w.Austrian)
			if err != nil {
				return err
			}
		}

		err := v.SetCursor(0, 0)
		if err != nil {
			return err
		}

		err = v.SetOrigin(0, 0)
		if err != nil {
			return err
		}
	}

	return nil
}