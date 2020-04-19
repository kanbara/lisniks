package main

import (
	"fmt"
	"github.com/kanbara/lisniks/pkg/dictionary"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"math/rand"
	"os"
	"time"
)

var (
	app  = kingpin.New("lisniks", "a reader for PolyGlot dictionaries")
	dictFile = app.Arg("dictionary", "the dictionary to open").Required().String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	dict := dictionary.Load(*dictFile)

	log.Infof("loaded dictionary from PolyGlot version %v, updated %v, word count %v",
		dict.Version, dict.LastUpdated, len(dict.Lexicon))
	log.Infof("%v - %v\n\n", dict.LanguageProperties.Name, dict.LanguageProperties.Version())

	// simply print out a few random words to see that our dictionary reading works as intended
	rand.Seed(time.Now().Unix())

	for i := 0; i <= 5; i++ {
		loc := rand.Intn(len(dict.Lexicon))
		word := dict.Lexicon[loc]

		fmt.Printf("%v (%v), type %v #%v\n\tdef: %v\n",
			word.Con, word.Local, word.Type, word.WordID, word.Def)
		fmt.Printf("------------------\n\n")
	}
}