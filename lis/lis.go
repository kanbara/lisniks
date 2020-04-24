package main

import (
	"github.com/jroimartin/gocui"
	"github.com/kanbara/lisniks/pkg/dictionary"
	"github.com/kanbara/lisniks/pkg/ui"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

var (
	app      = kingpin.New("lisniks", "a reader for PolyGlot dictionaries")
	dictFile = app.Arg("dictionary", "the dictionary to open").Required().String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	dict := dictionary.NewDictFromFile(*dictFile)
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicf("could not instantiate UI: %v", err)
	}

	defer g.Close()

	g.SetManager(ui.NewManager(dict))
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitfn); err != nil {
		log.Panicf("Could not handle keybinding: %v", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicf("MainLoop() errored: %v", err)
	}
}

func quitfn(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}
