package debug

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

// print Debug line to definition view
const debugViewName = "debug"

func GetOrCreateDebugView(g *gocui.Gui) (v *gocui.View) {
	if v, err := g.View(debugViewName); gocui.IsUnknownView(err) {
		maxX, maxY := g.Size()
		if v, err := g.SetView(debugViewName, 21, 10, maxX-1, maxY-6, 0); err != nil {

			v.Autoscroll = true
			v.BgColor = gocui.ColorRed
			v.FgColor = gocui.ColorWhite
			_, _ = g.SetViewOnTop(debugViewName)
			fmt.Fprintln(v, "DEBUG VIEW")

			return v
		}
	} else {
		return v
	}

	return nil
}

func MustRemoveDebugView(g *gocui.Gui) {
	_ = g.DeleteView(debugViewName)
}

func Clear(v *gocui.View) {
	if v == nil {
		return
	}

	v.Clear()

	if err := v.SetCursor(0, 0); err != nil {
		panic(err)
	}

	if err := v.SetWritePos(0, 10); err != nil {
		panic(err)
	}

	return
}

func Print(v *gocui.View, str string) {
	if _, err := fmt.Fprintln(v, str); err != nil {
		panic(err)
	}
}
