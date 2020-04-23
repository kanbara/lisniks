package language

import (
	"encoding/xml"
	"github.com/kanbara/lisniks/pkg/strings"
)

type AlphaOrderMap map[rune]int64

type Properties struct {
	Name       strings.Rawstring `xml:"langName"`
	Copyright  strings.Rawstring `xml:"langPropAuthorCopyright"`
	AlphaOrder AlphaOrderMap     `xml:"alphaOrder"`
}

// UnmarshalXML turns a tuple string like 2,4 into the appropriate struct {2,4}
// which can be used to get the resultant {Noun, Feminine} from a lookup
func (a *AlphaOrderMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	// we cast s to a slice of runes to get the proper character size
	// e.g. ǆ or č may take up more than one byte, which is the default
	// unit that makes up a string (it's a []byte effectively)

	runes := []rune(s)
	ao := make(AlphaOrderMap, len(runes))

	for i := 0; i < len(runes); i++ {
		ao[runes[i]] = int64(i)
	}



	// the unmarshaler calls the function on an empty instance of this type
	// and the data is stored in the xml.StartElement
	// this means we actually set the pointer of this type to the instance we create
	// instead of returning it
	*a = ao

	return nil
}

// Version gets the data from the Copyright field, as we assume that's where the version is stored
func (l Properties) Version() string {
	return l.Copyright.String()
}