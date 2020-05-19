package ui

import "github.com/awesome-gocui/gocui"


type ListView struct {
	View
	viewName string
	itemLen func() int
	itemSelected func() *int
	selected func(g *gocui.Gui, v *gocui.View) error
}

func (l *ListView) updatePosition(v *gocui.View, updown int) error {
	if l.itemLen() == 0 {
		return nil
	}

	// the cursorPos highlighted position in the viewSize e.g. which row selected
	cx, cy := v.Cursor()

	// we can't go above len words
	maxY := l.itemLen() - 1

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

	c, sel := l.calculateNewViewAndState(coords{
		cursorPos:   coordinates{cx, cy},
		originStart: coordinates{ox, oy},
		viewSize:    coordinates{vx, vy},
	}, updown, *l.itemSelected(), 0, maxY)

	err := v.SetCursor(c.cursorPos.x, c.cursorPos.y)
	if err != nil {
		return err
	}

	err = v.SetOrigin(c.originStart.x, c.originStart.y)
	if err != nil {
		l.Log.Debugf("starting cursor(%v,%v)\n",
			"starting origin(%v,%v)\n",
			"viewsize(%v,%v)\n",
			"updown %v\n",
			"sel %v\n",
			"0,%v\n",
			"output****\n",
			"%#v\n",
			"sel: %v", cx, cy, ox, oy, vx, vy, updown, *l.itemSelected(), maxY, c, sel)

		return err
	}


	// this looks weird, but it's just a function which returns a pointer, which we then
	// dereference to set sel to that value!

	*(l.itemSelected()) = sel
	l.UpdateViews()


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

func (l *ListView) calculateNewViewAndState(c coords, updown int,
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
		l.Log.Debugf("scrolling up out of frame")

		// we also need to check if we'll exceed the min entries
		// if the selected word would go below the min
		if c.originStart.y+updown < minY {
			l.Log.Debugf("scrolling past min entry")
			selected = 0
			out.originStart.y = 0
			out.cursorPos.y = 0
			break
		}

		out.originStart.y = c.originStart.y + updown
		selected = selected + updown

	case c.cursorPos.y+updown >= c.viewSize.y && maxY >= c.viewSize.y:
		// we are scrolling down out the frame
		l.Log.Debugf("scrolling down out of frame")

		if maxY-(c.originStart.y+updown) < c.viewSize.y-1 {
			l.Log.Debug("will have incomplete frame")

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
		l.Log.Debugf("scrolling past small list")
		out.cursorPos.y = maxY
		selected = maxY
	default: // we are inside the frame
		out.cursorPos.y = c.cursorPos.y + updown
		selected = selected + updown
	}

	return out, selected
}

func (l *ListView) nextItem(g *gocui.Gui, v *gocui.View) error {
	err := l.updatePosition(v, 1)
	if err != nil {
		return err
	}

	return nil
}

func (l *ListView) jump(v *gocui.View, ahead bool) error {
	var jump int
	_, vy := v.Size()
	// if there's not so many words, make a smaller jump
	// e.g. with word length less than one frame, so that we can
	// jump inside of a smaller list
	if vy > l.itemLen() {
		jump = l.itemLen() / 2
	} else {
		jump = vy
	}

	if !ahead {
		jump = -jump
	}

	err := l.updatePosition(v, jump)
	if err != nil {
		return err
	}

	return nil
}

func (l *ListView) nextItemJump(_ *gocui.Gui, v *gocui.View) error {
	err := l.jump(v, true)
	if err != nil {
		return err
	}

	return nil
}

func (l *ListView) prevItem(_ *gocui.Gui, v *gocui.View) error {
	err := l.updatePosition(v, -1)
	if err != nil {
		return err
	}

	return nil
}

func (l *ListView) prevItemJump(_ *gocui.Gui, v *gocui.View) error {
	err := l.jump(v, false)
	if err != nil {
		return err
	}

	return nil
}

func (l *ListView) SetKeybindings() error {
	if err := l.g.SetKeybinding(l.viewName, gocui.KeyEnter, gocui.ModNone, l.selected); err != nil {
		return err
	}

	if err := l.g.SetKeybinding(l.viewName, gocui.KeyArrowDown, gocui.ModNone, l.nextItem); err != nil {
		return err
	}

	// TODO make these configurable or at least a colemak option ;)
	if err := l.g.SetKeybinding(l.viewName, 'j', gocui.ModNone, l.nextItem); err != nil {
		return err
	}

	if err := l.g.SetKeybinding(l.viewName, gocui.KeyArrowUp, gocui.ModNone, l.prevItem); err != nil {
		return err
	}

	if err := l.g.SetKeybinding(l.viewName, 'k', gocui.ModNone, l.prevItem); err != nil {
		return err
	}

	if err := l.g.SetKeybinding(l.viewName, gocui.KeyCtrlF, gocui.ModNone, l.nextItemJump); err != nil {
		return err
	}

	if err := l.g.SetKeybinding(l.viewName, gocui.KeyCtrlB, gocui.ModNone, l.prevItemJump); err != nil {
		return err
	}

	return nil
}