package main

import (
	"fmt"
	"github.com/kanbara/lisniks/pkg/dictionary"
	"github.com/kanbara/lisniks/pkg/state"
	"github.com/kanbara/lisniks/pkg/ui"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
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

	vm := ui.NewViewManager(dict, s, logger)
	if err := vm.Run(); err != nil {
		logger.Panicf("error on main runloop: %v", err)
	}

}
