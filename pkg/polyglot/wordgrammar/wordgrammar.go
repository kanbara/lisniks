package wordgrammar

import (
	"encoding/xml"
	"github.com/kanbara/lisniks/pkg/polyglot/word"
	"strconv"
	gostrings "strings"
)

type Class []ClassNode

type Service struct {
	wgMap Map
}

type ApplyType []int64

type ClassNode struct {
	ID   int64  `xml:"wordGrammarClassID"`
	Name string `xml:"wordGrammarClassName"`
	// TODO this should be a slice
	ApplyTypes ApplyType    `xml:"wordGrammarApplyTypes"` // part of speech this applies to
	Values     []ClassValue `xml:"wordGrammarClassValuesCollection>wordGrammarClassValueNode"`
}

type ClassValue struct {
	ID   int64  `xml:"wordGrammarClassValueId"`
	Name string `xml:"wordGrammarClassValueName"`
}

type MapKey struct {
	word.Class
	applyType int64
}

type MapValue struct {
	ClassName string
	ValueName string
}

// Map holds all the WordClasses for O(1) lookup
// use the Class and Value which is an entry in the list of word grammars, like gender
// and the subkey like masculine
// e.g. 2 (gender)
//     / \
//    0   1
//   masc fem
//
// where the map has entries {2,0} -> masc, {2,1} -> fem
type Map map[MapKey]MapValue

func (s *Service) Get(applyType int64, class word.Class) *MapValue {
	val, ok := s.wgMap[MapKey{class, applyType}]
	if !ok {
		return nil
	}

	return &val
}

func NewWordGrammarService(classes Class) *Service {
	m := make(Map, len(classes))

	for _, c := range classes {
		for _, v := range c.Values {
			for _, a := range c.ApplyTypes {
				m[MapKey{
					Class:     word.Class{Class: c.ID, Value: v.ID},
					applyType: a,
				}] = MapValue{c.Name, v.Name}
			}
		}
	}

	s := Service{wgMap: m}
	return &s
}

// UnmarshalXML turns a tuple string like 2,4 into the appropriate struct {2,4}
// which can be used to get the resultant type for a given word
func (a *ApplyType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var list ApplyType
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	elems := gostrings.Split(s, ",")

	for _, e := range elems {
		i, err := strconv.ParseInt(e, 10, 64)
		if err != nil {
			return err
		}

		list = append(list, i)
	}

	*a = list

	return nil
}
