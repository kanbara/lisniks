package dictionary

import (
	"archive/zip"
	"encoding/xml"
	"github.com/c2h5oh/datasize"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

const dictXMLName string = "pgdictionary.xml"

func isZip(fileName string) bool {
	z, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Error(err)
		return false
	}

	zipMagic := [4]byte{0x50, 0x4B, 0x03, 0x04}   // PK..
	emptyZip := [4]byte{0x50, 0x4B, 0x05, 0x06}   // empty z
	spannedZip := [4]byte{0x50, 0x4B, 0x07, 0x08} // spanned z

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

func Load(filename string) Dictionary {
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

	dict := Dictionary{}
	// unmarshal the XML data into a Dictionary struct
	err = xml.Unmarshal(dictBytes, &dict)
	if err != nil {
		log.Fatalf("could not unmarshal bytes->dict: %v", err)
	}

	return dict
}
