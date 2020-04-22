package declension

type Service struct {
	declensions Declensions
}

type Declensions []Declension

type Declension struct {
	ID        int64  `xml:"declensionId"`
	Text      string `xml:"declensionText"`
	Notes     string `xml:"declensionNotes"`
	Template  int64  `xml:"declensionTemplate"`
	RelatedID int64  `xml:"declensionRelatedId"`
	Nodes     []Node `xml:"dimensionNode"`
}

type Node struct {
	ID   int64  `xml:"dimensionId"`
	Name string `xml:"dimensionName"`
}

func (s *Service) Get(id int64) Declension {
	return s.declensions[id]
}

func NewDeclensionService(declensions Declensions) *Service {
	s := Service{declensions: declensions}
	return &s
}
