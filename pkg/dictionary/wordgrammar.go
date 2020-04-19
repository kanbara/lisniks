package dictionary

type WordGrammarClass []WordGrammarClassNode

type WordGrammarClassNode struct {
	ID string `xml:"wordGrammarClassID"`
	Name string `xml:"wordGrammarClassName"`
	// TODO this should be a slice
	ApplyTypes string `xml:"wordGrammarApplyTypes"` // part of speech this applies to
	Values []WordGrammarClassValue `xml:"wordGrammarClassValuesCollection>wordGrammarClassValueNode"`
}

type WordGrammarClassValue struct {
	ID string `xml:"wordGrammarClassValueId"`
	Name string `xml:"wordGrammarClassValueName"`
}