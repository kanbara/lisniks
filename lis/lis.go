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
	app      = kingpin.New("lisniks", "a reader for PolyGlot dictionaries")
	dictFile = app.Arg("dictionary", "the dictionary to open").Required().String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	dict := dictionary.NewDictFromFile(*dictFile)

	log.Info(dict.Stats())
	log.Info(dict.LangAndVersion())
	// simply print out a few random words to see that our dictionary reading works as intended
	rand.Seed(time.Now().Unix())

	for i := 0; i <= 5; i++ {
		loc := rand.Intn(dict.Lexicon.Count())

		fmt.Print(dict.PrettyWordStringByLoc(loc))
		fmt.Print("------------------\n\n")
	}

	// temporary, while i i work on filtering out the extraneous classes
	fmt.Println(dict.PrettyWordStringByID(2550))
	fmt.Println(dict.PrettyWordStringByID(2554))


}
