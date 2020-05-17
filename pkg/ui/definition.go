package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"strings"
)

type DefinitionView struct {
	DefaultView
}

const spaceWidth = 1

func (d *DefinitionView) New(name string) error {

	maxX, maxY := d.g.Size()

	if v, err := d.g.SetView(name, 21, 10, maxX-1, maxY-6, 0); err != nil {
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
	if d.State.CurrentWord() != nil {
		//fmt.Printf(fmt.Sprintf("\n\n\n\n%q", d.state.CurrentWord().Def.String()))
		modified := wordWrap(d.State.CurrentWord().Def.String(), x-1)
				_, err := fmt.Fprintln(v, modified)
				if err != nil {
					return err
				}
	}

	return nil
}

func wordWrap(word string, lineWidth int) string {
	//shamelessely lifted this from the stupid internet and implemented it here
	spaceLeft := lineWidth
	var out string


	wordsSpace := strings.Split(word, " ")
	var words []string

	// go thru each word to check if we have newlines, ffs
	for _, w := range wordsSpace {
		tmp := ""
		for _, c := range w {
			// fml finding newlines
			// how can nothing like this exist srsly wtf
			// i just want to split a string on char C and
			// keep the C in the slice! (×﹏×)

			// if we have no newline, keep growing the word
			if c != '\n' {
				tmp += string(c)
			// if we have a newline, append the word and newline
			// and reset the word tmp thing
			} else {
				words = append(words, []string{tmp, "\n"}...)
				tmp = ""
			}
		}

		// if we get to the end with no newline, e.g. our word is ["this"]
		// write it out
		if tmp != "" {
			words = append(words, tmp)
		}
	}

	for i, w := range words {

		// there may already be \n's here
		if w == "\n" {
			out += "\n"
			spaceLeft = lineWidth
			continue
		}

		// word+space is greater than space left on line
		// so we add a newline
		if len([]rune(w)) + spaceWidth > spaceLeft {
			out += "\n"

			// of course the space is now subtracted from the word
			// that's at the start of the next line
			spaceLeft = lineWidth - len([]rune(w))
		} else {
			// just subtract the word and space width from the remaining width on the line

			// don't put a space at first
			//also don't put a space if we just had a newline
			if i != 0 && words[i-1] != "\n" {
					out += " "
			}
			spaceLeft = spaceLeft - (len([]rune(w)) + spaceWidth)
		}

		// add the word to our string
		out += w
	}

	return out
}
