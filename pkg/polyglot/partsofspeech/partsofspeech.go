package partsofspeech

import (
	"github.com/kanbara/lisniks/pkg/strings"
)

type Service struct {
	posMap Map
	reversePosMap ReverseMap
}

type PartsOfSpeech []Part

type Part struct {
	ID    int64             `xml:"classId"`
	Name  string            `xml:"className"`
	Notes strings.Rawstring `xml:"classNotes"`
}

type Map map[int64]string
type ReverseMap map[string]int64

func (s *Service) GetByID(id int64) string {
	return s.posMap[id]
}

func (s *Service) GetByName(str string) int64 {
	return s.reversePosMap[str]
}

func NewPartsOfSpeechService(pos PartsOfSpeech) *Service {

	m := make(Map, len(pos))
	rm := make(ReverseMap, len(pos))


	for _, p := range pos {
		m[p.ID] = p.Name
		rm[p.Name] = p.ID
	}

	s := Service{posMap: m, reversePosMap: rm}
	return &s
}
