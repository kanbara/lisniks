package dictionary

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type Word struct {
	WordID int64 `xml:"wordId"`
	Local Rawstring `xml:"localWord"`
	Con Rawstring `xml:"conWord"`
	Type int64 `xml:"wordTypeId"`
	Def Rawstring `xml:"definition"`
	Classes []WordClass `xml:"wordClassCollection>wordClassification"`
}

// WordClass allows us to add a custom unmarshaler to get the values out in a way
// that makes is easy to lookup later
type WordClass struct {
	Class int64
	Value int64
}

// UnmarshalXML turns a tuple string like 2,4 into the appropriate struct {2,4}
// which can be used to get the resultant {Noun, Feminine} from a lookup
func (wc *WordClass) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	elems := strings.Split(s, ",")

	// string -> int64
	class, err := strconv.ParseInt(elems[0], 10, 64)
	if err != nil {
		return err
	}

	// string -> int64
	val, err := strconv.ParseInt(elems[1], 10, 64)
	if err != nil {
		return err
	}

	// the unmarshaler calls the function on an empty instance of this type
	// and the data is stored in the xml.StartElement
	// this means we actually set the pointer of this type to the instance we create
	// instead of returning it
	*wc = WordClass{class, val}

	return nil
}