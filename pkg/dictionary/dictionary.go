package dictionary

import (
	"archive/zip"
	"encoding/xml"
	"github.com/c2h5oh/datasize"
	"github.com/kanbara/lisniks/pkg/lexicon"
	"github.com/kanbara/lisniks/pkg/partsofspeech"
	"github.com/kanbara/lisniks/pkg/wordgrammar"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

const dictXMLName string = "pgdictionary.xml"

type Dictionary struct {
	file *File

	PartsOfSpeech *partsofspeech.Service
	WordGrammar   *wordgrammar.Service
	Lexicon       *lexicon.Service
}

// NewDictFromFile will load the internal XML dictionary from a PolyGlot ZIP to a DictionaryFile struct in memory
func NewDictFromFile(filename string) *Dictionary {
	dictFile := mustGetDictFileFromXML(filename)
	dict := mustUnmarshalXML(dictFile)

	d := Dictionary{file: &dict}
	d.PartsOfSpeech = partsofspeech.NewPartsOfSpeechService(dict.PartsOfSpeech)
	d.WordGrammar = wordgrammar.NewWordGrammarService(dict.WordGrammarClasses)
	d.Lexicon = lexicon.NewLexiconService(dict.Lexicon)

	return &d
}

// mustUnmarshalXML takes a file from the zip and unmarshals it to a go dictionary.File structure
// or exists on error
func mustUnmarshalXML(dictFile *zip.File) File {
	// get a reader for the dictionary file
	dictRC, err := dictFile.Open()
	if err != nil {
		log.Fatalf("could not open dictionary xml: %v", err)
	}

	// read the contents of the file
	dictBytes, err := ioutil.ReadAll(dictRC)
	if err != nil {
		log.Fatalf("could not read dictionary xml: %v", err)
	}

	dict := File{}
	// unmarshal the XML data into a DictionaryFile struct
	err = xml.Unmarshal(dictBytes, &dict)
	if err != nil {
		log.Fatalf("could not unmarshal bytes->dict: %v", err)
	}
	return dict
}

// mustGetDictFileFromXML takes a filename and returns a zip or exits on error
func mustGetDictFileFromXML(filename string) *zip.File {
	ok := isZip(filename)
	if !ok {
		log.Fatalf("`%v` is not a zip, bailing", filename)
	}

	log.Infof("`%v` is zip; extracting dict XML", filename)

	rc, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal(err)
	}

	var dictFile *zip.File
	files := rc.File
	// try to find the dictionary XML file inside the zip
	for _, f := range files {
		if strings.ToLower(f.Name) == dictXMLName {
			dictFile = f
			break
		}
	}

	if dictFile == nil {
		log.Fatalf("could not find %v in dictionary zip", dictXMLName)
		os.Exit(1) // this code would never be reached, just here to calm the linter
	}

	log.Infof("found dictionary xml with %v bytes",
		(datasize.ByteSize(dictFile.UncompressedSize64) * datasize.B).HumanReadable())

	return dictFile
}

// isZip determines if a file is reasonably a zip or not.
// we don't really need this, because i am sure zip.OpenReader will error otherwise
// but it can't hurt i suppose
func isZip(fileName string) bool {
	z, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Error(err)
		return false
	}

	// these are the 3 magic numbers for different zip types
	zipMagic := [4]byte{0x50, 0x4B, 0x03, 0x04}   // PK..
	emptyZip := [4]byte{0x50, 0x4B, 0x05, 0x06}   // empty zip
	spannedZip := [4]byte{0x50, 0x4B, 0x07, 0x08} // spanned zip

	// get the 4 byte array to compare with the header
	zipHdr := [4]byte{z[0], z[1], z[2], z[3]}

	switch {
	case zipHdr == zipMagic:
		return true
	case zipHdr == emptyZip:
		log.Error("empty zip found")
		return false
	case zipHdr == spannedZip:
		log.Error("spanned zip found")
		return false
	default:
		log.Error("unknown file format found")
		return false
	}
}
