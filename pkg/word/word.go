package word

import (
	"encoding/xml"
	"github.com/kanbara/lisniks/pkg/strings"
	"strconv"
	gostrings "strings"
)

type Word struct {
	WordID  int64             `xml:"wordId"`
	Local   strings.Rawstring `xml:"localWord"`
	Con     strings.Rawstring `xml:"conWord"`
	Type    int64             `xml:"wordTypeId"`
	Def     strings.Rawstring `xml:"definition"`
	Classes []Class           `xml:"wordClassCollection>wordClassification"`
}

// WordClass allows us to add a custom unmarshaler to get the values out in a way
// that makes is easy to lookup later
type Class struct {
	Class int64
	Value int64
}

// UnmarshalXML turns a tuple string like 2,4 into the appropriate struct {2,4}
// which can be used to get the resultant {Noun, Feminine} from a lookup
func (wc *Class) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	elems := gostrings.Split(s, ",")

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
	*wc = Class{class, val}

	return nil
}
