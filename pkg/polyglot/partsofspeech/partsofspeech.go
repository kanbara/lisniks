package partsofspeech

import (
	"github.com/kanbara/lisniks/pkg/strings"
	"sort"
)

type Service struct {
	posMap     Map
	nameLookup NameLookup
}

type PartsOfSpeech []Part

type Part struct {
	ID    int64             `xml:"classId"`
	Name  string            `xml:"className"`
	Notes strings.Rawstring `xml:"classNotes"`
}

type Map map[int64]string
type ReverseMap map[string]int64

type NameLookup []Part

// sort functions
// needed to make sort.Sort work on Service
// the only weird thing is less,
// because as opposed to normal strings, we need to sort on the Alpha Word Order
// from the dictionary, e.g. -AaĀāæ...Žž
// to ensure we order words the same as PolyGlot does

func (n NameLookup) Len() int {
	return len(n)
}

func (n NameLookup) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n NameLookup) Less(i, j int) bool {

	return n[i].Name < n[j].Name
}

func (s *Service) GetNameToIDs() NameLookup {
	return s.nameLookup
}

func (s *Service) GetByID(id int64) string {
	return s.posMap[id]
}

func (s *Service) GetByName(str string) int64 {
	for _, n := range s.nameLookup {
		if n.Name == str {
			return n.ID
		}
	}

	return 0
}

func NewPartsOfSpeechService(pos PartsOfSpeech) *Service {

	m := make(Map, len(pos))
	var n NameLookup

	for _, p := range pos {
		m[p.ID] = p.Name
		n = append(n, Part{Name: p.Name, ID: p.ID})
	}

	sort.Sort(n)
	s := Service{posMap: m, nameLookup: n}
	return &s
}
