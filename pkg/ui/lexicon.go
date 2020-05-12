package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	log "github.com/sirupsen/logrus"
)

type LexiconView struct {
	View
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
		v.Title = fmt.Sprintf("%v %v/%v",name, len(l.state.Words), l.dict.Lexicon.Len())
		v.Frame = true
		v.Highlight = true
		v.SelFgColor = gocui.ColorGreen
		v.FgColor = gocui.ColorWhite

		for _, w := range l.state.Words {
			_, err := fmt.Fprintln(v, w.Austrian)
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

		err = v.SetHighlight(l.state.SelectedWord, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LexiconView) Update(v *gocui.View) error {
	v.Clear()
	l.state.SelectedWord = 0
	v.Title = fmt.Sprintf("%v %v/%v", lexView, len(l.state.Words), l.dict.Lexicon.Len())

	if len(l.state.Words) > 0 {
		for _, w := range l.state.Words {
			_, err := fmt.Fprintln(v, w.Austrian)
			if err != nil {
				return err
			}
		}

		err := v.SetHighlight(0, true)
		if err != nil {
			return err
		}

		err = v.SetCursor(0, 0)
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

func (l *LexiconView) SetKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(lexView, gocui.KeyArrowDown, gocui.ModNone, l.NextWord); err != nil {
		return err
	}

	if err := g.SetKeybinding(lexView, gocui.KeyArrowUp, gocui.ModNone, l.prevWord); err != nil {
		return err
	}

	if err := g.SetKeybinding(lexView, gocui.KeyCtrlF, gocui.ModNone, l.nextWordJump); err != nil {
		return err
	}

	if err := g.SetKeybinding(lexView, gocui.KeyCtrlB, gocui.ModNone, l.prevWordJump); err != nil {
		return err
	}

	return nil
}

func (l *LexiconView) updateWord(g *gocui.Gui, v *gocui.View, updown int) error {
	if len(l.state.Words) == 0 {
		return nil
	}

	// the cursorPos highlighted position in the viewSize e.g. which row selected
	cx, cy := v.Cursor()

	// we can't go above len words
	maxY := len(l.state.Words) - 1

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
	err := v.SetHighlight(l.state.SelectedWord, false)
	if err != nil {
		return err
	}

	c, sel := calculateNewViewAndState(coords{
		cursorPos:   coordinates{cx, cy},
		originStart: coordinates{ox, oy},
		viewSize:    coordinates{vx, vy},
	}, updown, l.state.SelectedWord, 0, maxY)

	err = v.SetCursor(c.cursorPos.x, c.cursorPos.y)
	if err != nil {
		return err
	}

	err = v.SetOrigin(c.originStart.x, c.originStart.y)
	if err != nil {
		// todo this info is super useful, let's wrap the error
		//panic(fmt.Sprintf("starting cursor(%v,%v)\n" +
		//	"starting origin(%v,%v)\n" +
		//	"viewsize(%v,%v)\n" +
		//	"updown %v\n" +
		//	"sel %v\n" +
		//	"0,%v\n" +
		//	"output****\n" +
		//	"%#v\n" +
		//	"sel: %v", cx,cy,ox,oy,vx,vy,updown,m.state.SelectedWord,maxY,c,sel))
		return err
	}

	l.state.SelectedWord = sel

	// highlight current word
	err = v.SetHighlight(l.state.SelectedWord, true)
	if err != nil {
		return err
	}

	g.Update(func(g *gocui.Gui) error {
		for _, viewName := range l.viewsToUpdate {
			if v, err := g.View(viewName); err != nil {
				return err
			} else {
				if err := l.views[viewName].Update(v); err != nil {
					return err
				}
			}
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
	case c.cursorPos.y+updown < 0: // we are scrolling up out of the frame
		log.Debugf("scrolling up out of frame")
		log.Debugf("%v + %v < %v", c.cursorPos.y, updown, 0)

		// we also need to check if we'll exceed the min entries
		// if the selected word would go below the min
		if c.originStart.y+updown < minY {
			log.Debugf("scrolling past min entry")
			selected = 0
			out.originStart.y = 0
			out.cursorPos.y = 0
			break
		}

		out.originStart.y = c.originStart.y + updown
		selected = selected + updown

	case c.cursorPos.y+updown >= c.viewSize.y && maxY >= c.viewSize.y:
		// we are scrolling down out the frame
		log.Debugf("scrolling down out of frame")
		log.Debugf("%v + %v >= %v", c.cursorPos.y, updown, c.viewSize.y)

		if maxY-(c.originStart.y+updown) < c.viewSize.y-1 {
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
	case c.cursorPos.y+updown > maxY: // scrolling past a small list
		log.Debugf("scrolling past small list")
		log.Debugf("%v + %v > %v", c.cursorPos.y, updown, maxY)
		out.cursorPos.y = maxY
		selected = maxY
	default: // we are inside the frame
		log.Debugf("default: cursor   %v + %v", c.cursorPos.y, updown)
		log.Debugf("default: selected %v + %v", selected, updown)

		out.cursorPos.y = c.cursorPos.y + updown
		selected = selected + updown
	}

	return out, selected
}

func (l *LexiconView) NextWord(g *gocui.Gui, v *gocui.View) error {
	err := l.updateWord(g, v, 1)
	if err != nil {
		return err
	}

	return nil
}

func (l *LexiconView) jump(g *gocui.Gui, v *gocui.View, ahead bool) error {
	var jump int
	_, vy := v.Size()
	// if there's not so many words, make a smaller jump
	// e.g. with word length less than one frame, so that we can
	// jump inside of a smaller list
	if vy > len(l.state.Words) {
		jump = len(l.state.Words) / 2
	} else {
		jump = vy
	}

	if !ahead {
		jump = -jump
	}

	err := l.updateWord(g, v, jump)
	if err != nil {
		return err
	}

	return nil
}

func (l *LexiconView) nextWordJump(g *gocui.Gui, v *gocui.View) error {
	err := l.jump(g, v, true)
	if err != nil {
		return err
	}

	return nil
}

func (l *LexiconView) prevWord(g *gocui.Gui, v *gocui.View) error {
	err := l.updateWord(g, v, -1)
	if err != nil {
		return err
	}

	return nil
}

func (l *LexiconView) prevWordJump(g *gocui.Gui, v *gocui.View) error {
	err := l.jump(g, v, false)
	if err != nil {
		return err
	}

	return nil
}
