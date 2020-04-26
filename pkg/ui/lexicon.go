package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	log "github.com/sirupsen/logrus"
)

func (m *Manager) NewLexiconView(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView(lexView, 0, 6, 20, maxY-1, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		// TODO i wanted to set the title to the current word
		// but unfortunately it didn't work because it printed
		// sth. like -jeř-a because of the two byte width unicode
		//
		// i saw an issue on the thing about unicode, maybe there's a fix
		// that can be done.
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
	// the cursorPos highlighted position in the viewSize e.g. which row selected
	cx, cy := v.Cursor()

	// we can't go above len words
	maxY := len(m.state.Words)-1

	// the position of the buffer in the viewSize
	ox, oy := v.Origin()

	// this is the max y of the viewSize e.g. how many lines of words are currently in the viewSize
	vx, vy := v.Size()

	// if we will exceed the frame of the lexicon
	// we need to take care to move the originStart of the buffer instead of the
	// cursorPos
	// e.g.
	//
	// ------
	// a
	// b
	// c <--
	// ------
	// d
	// e
	//
	// when we are at position c, we want to move originStart down one
	// instead of moving the cursorPos as we normally would
	// cursorPos == highlight in this case

	// turn off highlight for previous words
	err := v.SetHighlight(m.state.SelectedWord, false)
	if err != nil {
		return err
	}

	c, sel := calculateNewViewAndState(coords{
		cursorPos:   coordinates{cx,cy},
		originStart: coordinates{ox, oy},
		viewSize:    coordinates{vx, vy},
	}, updown, m.state.SelectedWord, 0, maxY)

	err = v.SetCursor(c.cursorPos.x, c.cursorPos.y)
	if err != nil {
		return err
	}

	err = v.SetOrigin(c.originStart.x, c.originStart.y)
	if err != nil {
		return err
	}

	m.state.SelectedWord = sel

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

		v, err = g.View(currentWordView)
		if err != nil {
			return err
		}

		err = m.updateCurrentWordView(v)
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

type coordinates struct {
	x, y int
}

type coords struct {
	cursorPos   coordinates
	originStart coordinates
	viewSize    coordinates
}

func calculateNewViewAndState(c coords, updown int,
	selected int, minY int, maxY int) (coords, int) {

	out := coords{
		cursorPos: coordinates{
			x: c.cursorPos.x,
			y: c.cursorPos.y,
		},
		originStart: coordinates{
			x: c.originStart.x,
			y: c.originStart.y,
		},
	}

	switch {
	case c.cursorPos.y + updown < 0: // we are scrolling up out of the frame
		log.Debugf("scrolling up out of frame")
		log.Debugf("%v + %v < %v", c.cursorPos.y, updown, 0)

		// we also need to check if we'll exceed the min entries
		// if the selected word would go below the min
		if c.originStart.y + updown < minY {
			log.Debugf("scrolling past min entry")
			selected = 0
			out.originStart.y = 0
			out.cursorPos.y = 0
			break
		}

		out.originStart.y = c.originStart.y + updown
		selected = selected + updown
	case c.cursorPos.y + updown >= c.viewSize.y: // we are scrolling down out the frame
		log.Debugf("scrolling down out of frame")
		log.Debugf("%v + %v >= %v", c.cursorPos.y, updown, c.viewSize.y)

		if maxY - (c.originStart.y + updown) < c.viewSize.y - 1 {
			log.Debug("will have incomplete frame")
			log.Debugf("%v + %v < %v", c.originStart.y, updown, c.viewSize.y)

			// i think this gets a + 1 because the maxY and viewsize
			// are ... yeah idk, maybe one is one based and the other not?
			// TODO
			// also the cursorpos gets a -1 because viewsize
			// is how many items fit in the view, so if we fit 30
			// we need to be at pos 29
			out.originStart.y = maxY - c.viewSize.y + 1
			out.cursorPos.y = c.viewSize.y - 1
			selected = maxY
			break
		}

		out.originStart.y = c.originStart.y + updown
		selected = selected + updown
	case c.cursorPos.y + updown >= maxY: // scrolling past a small list
		log.Debugf("scrolling past small list")
		log.Debugf("%v + %v > %v", c.cursorPos.y, updown, maxY)
		out.cursorPos.y = maxY
		selected = maxY
	default: // we are inside the frame
		log.Debugf("default: cursor   %v + %v" ,c.cursorPos.y, updown)
		log.Debugf("default: selected %v + %v" ,selected, updown)

		out.cursorPos.y = c.cursorPos.y + updown
		selected = selected + updown
	}

	return out, selected
}

func (m *Manager) nextWord(g *gocui.Gui, v *gocui.View) error {
	err := m.updateWord(g, v, 1)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) nextWordJump(g *gocui.Gui, v *gocui.View) error {
	// get reasonable jump size based on word list
	onepct := len(m.state.Words) / 100
	err := m.updateWord(g, v, onepct)
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
	// get reasonable jump size based on word list
	onepct := len(m.state.Words) / 100
	err := m.updateWord(g, v, -onepct)
	if err != nil {
		return err
	}

	return nil
}