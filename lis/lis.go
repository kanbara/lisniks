package main

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/kanbara/lisniks/pkg/dictionary"
	"github.com/kanbara/lisniks/pkg/state"
	"github.com/kanbara/lisniks/pkg/ui"
	 "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
)

var (
	app      = kingpin.New("lisniks", "a reader for PolyGlot dictionaries")
	dictFile = app.Arg("dictionary", "the dictionary to open").Required().String()
	debug    = app.Flag("debug", "debug").Short('d').Bool()

	// Version gets injected with -ldflags
	Version string
	// BuildTime gets injected with -ldflags
	BuildTime string
)

func newLogger(debug bool) *logrus.Logger {
	logger := logrus.New()
	if debug {
		logger.SetLevel(logrus.DebugLevel)
		file, _ := os.OpenFile("debug.log",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		logger.SetOutput(file)
		return logger
	}

	logger.SetOutput(ioutil.Discard)
	return logger
}

func main() {
	app.Version(fmt.Sprintf("%v, build time: %v", Version, BuildTime))
	kingpin.MustParse(app.Parse(os.Args[1:]))

	logger := newLogger(*debug)
	dict := dictionary.NewDictFromFile(*dictFile, logger)
	s := state.NewState(Version, dict)

	g, err := gocui.NewGui(gocui.Output256, false)
	if err != nil {
		log.Panicf("could not instantiate UI: %v", err)
	}

	defer g.Close()

	m := ui.NewManager(dict, s, logger)
	g.SetManager(m)

	err = m.SetGlobalKeybindings(g)
	if err != nil {
		log.Panicf("Could not set keybinding: %v", err)
	}

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.SelFrameColor = gocui.ColorGreen

	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		// debug stuff if we crash sometimes
		// not sure if useful but oh well
		view := g.CurrentView()
		ox, oy := view.Origin()
		vx, vy := view.Size()
		cx, cy := view.Cursor()
		cur, err := view.Line(cy)
		p := fmt.Sprintf("%v\nselected: %v\nview: %v\nview origin: %v,%v\n"+
			"view size: %v, %v\nview cursor: %v,%v\nlexicion list: %v\nbuf: `%v`",
			err, s.SelectedWord, view.Name(), ox, oy, vx, vy, cx, cy, len(s.Words), cur)
		panic(p)
	}
}
