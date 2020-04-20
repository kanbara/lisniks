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

// TODO you know the drill, move this somewhere better
// Also apparently we need to filter out the class that doesn't exist e.g.
// hassa is noun which should have gender and declensino, not CLASS4
// there should be type hints on the TypeID of the word grammar class
// a la ApplyTypes.
//
// yay more map lookups
func ClassesString(classes []dictionary.WordClass, wgMap dictionary.WGMap) string {
	var out string
	for _, c := range classes {
		val := wgMap[c]
		out += fmt.Sprintf("%v\n", val.ValueName)
	}

	return out
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	dict := dictionary.Load(*dictFile)

	// TODO there should be a root object created which handles these maps
	p := dictionary.MakePOSMap(dict.PartsOfSpeech)
	wg := dictionary.MakeWGMap(dict.WordGrammarClasses)

	log.Infof("loaded dictionary from PolyGlot version %v, updated %v, word count %v",
		dict.Version, dict.LastUpdated, len(dict.Lexicon))
	log.Infof("%v - %v\n\n", dict.LanguageProperties.Name, dict.LanguageProperties.Version())

	// simply print out a few random words to see that our dictionary reading works as intended
	rand.Seed(time.Now().Unix())

	for i := 0; i <= 5; i++ {
		loc := rand.Intn(len(dict.Lexicon))
		word := dict.Lexicon[loc]

		// TODO handler object should have a Word(WordID)
		// TODO which looks up the word in the lexicon and does this print
		fmt.Printf("%v (%v) [%v] #%v\n%v\n\tdef: %v\n",
			word.Con, word.Local, p[word.Type], word.WordID, ClassesString(word.Classes, wg), word.Def)
		fmt.Printf("------------------\n\n")
	}

	// temporary, while i i work on filtering out the extraneous classes
	word := dict.Lexicon.GetByID(2550)

	fmt.Printf("%v (%v) [%v] #%v\n%v\n\tdef: %v\n",
		word.Con, word.Local, p[word.Type], word.WordID, ClassesString(word.Classes, wg), word.Def)


	fmt.Printf("\n\n")
	word = dict.Lexicon.GetByID(2554)

	fmt.Printf("%v (%v) [%v] #%v\n%v\n\tdef: %v\n",
		word.Con, word.Local, p[word.Type], word.WordID, ClassesString(word.Classes, wg), word.Def)
}