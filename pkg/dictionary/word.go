package dictionary

type Word struct {
	WordID int64 `xml:"wordId"`
	Local Rawstring `xml:"localWord"`
	Con Rawstring `xml:"conWord"`
	Type int64 `xml:"wordTypeId"`
	Def Rawstring `xml:"definition"`
}