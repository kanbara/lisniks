package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"strings"
)

type DefinitionView NoBindingsView

const spaceWidth = 1

func (d *DefinitionView) New(g *gocui.Gui, name string) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView(name, 21, 10, maxX-1, maxY-6, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Title = name
		v.Wrap = true
		err := d.Update(v)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO at some point if we allow editing here,
// we have to save the word-break state as a `view` state and not just
// modify the word!
// esp because we get rid of all that weird html and stuff
func (d *DefinitionView) Update(v *gocui.View) error {
	v.Clear()

	x, _ := v.Size()
	if d.state.CurrentWord() != nil {

		modified := breakWordBoundaries(d.state.CurrentWord().Def.String(), x-2)
		_, err := fmt.Fprintln(v, modified)
		if err != nil {
			return err
		}
	}

	return nil
}

func breakWordBoundaries(word string, lineWidth int) string {
	//shamelessely lifted this from the stupid internet and implemented it here
	spaceLeft := lineWidth
	out := make([]string, 0, len(word))

	words := strings.Split(word, " ")
	for _, w := range words {
		if len([]rune(w)) + spaceWidth > spaceLeft {
			out = append(out, "\n")
			out = append(out, w)
			spaceLeft = lineWidth - len([]rune(w))
		} else {
			out = append(out, w)
			spaceLeft = spaceLeft - (len([]rune(w)) + spaceWidth)
		}
	}

	return strings.Join(out, " ")
}
