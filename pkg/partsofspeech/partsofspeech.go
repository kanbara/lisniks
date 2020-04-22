package partsofspeech

import (
	"github.com/kanbara/lisniks/pkg/strings"
)

type Service struct {
	posMap Map
}

type PartsOfSpeech []Part

type Part struct {
	ID    int64             `xml:"classId"`
	Name  string            `xml:"className"`
	Notes strings.Rawstring `xml:"classNotes"`
}

type Map map[int64]string

func (s *Service) Get(id int64) string {
	return s.posMap[id]
}

func NewPartsOfSpeechService(pos PartsOfSpeech) *Service {

	m := make(Map, len(pos))

	for _, p := range pos {
		m[p.ID] = p.Name
	}

	s := Service{posMap: m}
	return &s
}
