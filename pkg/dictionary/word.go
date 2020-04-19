package dictionary

type WordType string

type Word struct {
	WordID int64 `xml:"wordId"`
	Local string `xml:"localWord"`
	Con string `xml:"conWord"`
	Type WordType `xml:"wordTypeId"`
	Def Definition `xml:"definition"`
}