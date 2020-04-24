package main

import (
	"github.com/awesome-gocui/gocui"
	"github.com/kanbara/lisniks/pkg/dictionary"
	"github.com/kanbara/lisniks/pkg/state"
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
	s := state.State{
		Words: dict.Lexicon.Words(),
	}

	g, err := gocui.NewGui(gocui.Output256, false)
	if err != nil {
		log.Panicf("could not instantiate UI: %v", err)
	}

	defer g.Close()

	m := ui.NewManager(dict, &s)
	g.SetManager(m)

	err = m.SetKeybindings(g)
	if err != nil {
		log.Panicf("Could not set keybinding: %v", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicf("MainLoop() errored: %v", err)
	}
}
