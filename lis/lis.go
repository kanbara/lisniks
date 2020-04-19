package main

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/c2h5oh/datasize"
	"github.com/kanbara/lisniks/pkg/dictionary"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

const dictXMLName string = "pgdictionary.xml"

var (
	app  = kingpin.New("lisniks", "a reader for PolyGlot dictionaries")
	dict = app.Arg("dictionary", "the dictionary to open").Required().String()
)

func isZip(fileName string) error {
	zip, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("could not read dictionary: %w", err)
	}

	zipMagic := [4]byte{0x50, 0x4B, 0x03, 0x04}   // PK..
	emptyZip := [4]byte{0x50, 0x4B, 0x05, 0x06}   // empty zip
	spannedZip := [4]byte{0x50, 0x4B, 0x07, 0x08} // spanned zip

	zipHdr := [4]byte{zip[0], zip[1], zip[2], zip[3]}

	switch {
	case zipHdr == zipMagic:
		return nil
	case zipHdr == emptyZip:
		return errors.New("empty zip found")
	case zipHdr == spannedZip:
		return errors.New("spanned zip found")
	default:
		return errors.New("unknown file format found")
	}
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	err := isZip(*dict)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("%v is zip; extracting dict XML", *dict)

	rc, err := zip.OpenReader(*dict)
	if err != nil {
		log.Fatal(err)
	}

	var dictFile *zip.File
	files := rc.File
	for _, f := range files {
		if strings.ToLower(f.Name) == dictXMLName {
			dictFile = f
			break
		}
	}

	if dictFile == nil {
		log.Fatalf("could not find %v in dictionary zip", dictXMLName)
		os.Exit(1) // dummy call to please the linter, as Fatalf does this already
	}

	log.Infof("found dictionary xml with %v bytes",
		(datasize.ByteSize(dictFile.UncompressedSize64) * datasize.B).HumanReadable())

	dictRC, err := dictFile.Open()
	if err != nil {
		log.Fatalf("could not open dictionary xml: %v", err)
	}

	dictBytes, err := ioutil.ReadAll(dictRC)
	if err != nil {
		log.Fatalf("could not read dictionary xml: %v", err)
	}

	dict := dictionary.Dictionary{}
	err = xml.Unmarshal(dictBytes, &dict)
	if err != nil {
		log.Fatalf("could not unmarshal bytes->dict: %v", err)
	}

	log.Infof("loaded dictionary from PolyGlot version %v, updated %v, word count %v",
		dict.Version, dict.LastUpdated, len(dict.Lexicon))
	log.Infof("%v - %v", dict.LanguageProperties.Name, dict.LanguageProperties.Version())

	rand.Seed(time.Now().Unix())

	for i := 0; i <= 5; i++ {
		loc := rand.Intn(len(dict.Lexicon))
		word := dict.Lexicon[loc]

		fmt.Printf("%v (%v), type %v\n\tdef: %v\n\n",
			word.Con, word.Local, word.Type, word.Def)
	}
}
