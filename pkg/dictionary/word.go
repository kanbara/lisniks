package dictionary

type WordType string

type Word struct {
	WordID int64 `xml:"wordId"`
	Local Rawstring `xml:"localWord"`
	Con Rawstring `xml:"conWord"`
	Type WordType `xml:"wordTypeId"`
	Def Rawstring `xml:"definition"`
}