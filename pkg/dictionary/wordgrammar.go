package dictionary

type WordGrammarClass []WordGrammarClassNode

type WordGrammarClassNode struct {
	ID int64 `xml:"wordGrammarClassID"`
	Name string `xml:"wordGrammarClassName"`
	// TODO this should be a slice
	ApplyTypes string `xml:"wordGrammarApplyTypes"` // part of speech this applies to
	Values []WordGrammarClassValue `xml:"wordGrammarClassValuesCollection>wordGrammarClassValueNode"`
}

type WordGrammarClassValue struct {
	ID int64 `xml:"wordGrammarClassValueId"`
	Name string `xml:"wordGrammarClassValueName"`
}

// WordClassValue holds the Class name and Value name in case we want to introspect them both later
type WordClassValue struct {
	ClassName string
	ValueName string
}

// WGMap holds all the WordClasses for O(1) lookup
type WGMap map[WordClass]WordClassValue

func MakeWGMap(classes WordGrammarClass) WGMap {
	m := make(WGMap, len(classes))

	for _, c := range classes {
		for _, v := range c.Values {
			m[WordClass{c.ID, v.ID}] = WordClassValue{c.Name, v.Name}
		}
	}

	return m
}