package wordgrammar

import (
	"github.com/kanbara/lisniks/pkg/word"
)

type Class []ClassNode

type Service struct {
	wgMap Map
}

type ClassNode struct {
	ID   int64  `xml:"wordGrammarClassID"`
	Name string `xml:"wordGrammarClassName"`
	// TODO this should be a slice
	ApplyTypes string       `xml:"wordGrammarApplyTypes"` // part of speech this applies to
	Values     []ClassValue `xml:"wordGrammarClassValuesCollection>wordGrammarClassValueNode"`
}

type ClassValue struct {
	ID   int64  `xml:"wordGrammarClassValueId"`
	Name string `xml:"wordGrammarClassValueName"`
}

// WordClassValue holds the Class name and Value name in case we want to introspect them both later
type WordClassValue struct {
	ClassName string
	ValueName string
}

// WGMap holds all the WordClasses for O(1) lookup
type Map map[word.Class]WordClassValue

func (s *Service) Get(class word.Class) WordClassValue {
	return s.wgMap[class]
}

func NewWordGrammarService(classes Class) *Service {
	m := make(Map, len(classes))

	for _, c := range classes {
		for _, v := range c.Values {
			m[word.Class{Class: c.ID, Value: v.ID}] =
				WordClassValue{c.Name, v.Name}
		}
	}

	s := Service{wgMap: m}
	return &s
}
