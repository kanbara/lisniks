package state

import (
	"github.com/kanbara/lisniks/pkg/lexicon"
)

type SearchData struct {
	Type    lexicon.SearchType
	Pattern lexicon.SearchPattern
	String  string
}

type SearchQueue struct {
	data []SearchData
	max int
}

func NewSearchQueue(max int) SearchQueue {
	d := make([]SearchData, 0, max)
	return SearchQueue{max: max, data: d}
}

func (s *SearchQueue) RemoveOthers(sd SearchData) {
	n := 0
	nd := make([]SearchData, len(s.data)) // at most we'll have len(String.data) points
	for _, x := range s.data {
		if x.String != sd.String || // string compare first, as that'String more likely to short circuit
			x.Type != sd.Type ||
			x.Pattern != sd.Pattern {
			nd[n] = x // if x does not match sd, we keep it
			n++
		}
	}

	s.data = nd[:n]
}

func (s *SearchQueue) Enqueue(sd SearchData) {

	s.RemoveOthers(sd)

	if len(s.data)+1 > s.max {
		_ = s.Dequeue()
	}

	s.data = append(s.data, sd)
}

func (s *SearchQueue) Dequeue() *SearchData {
	if len(s.data) == 0 {
		return nil
	}

	n := SearchQueue{
		max: s.max,
		data: s.data[1:len(s.data)],
	}

	popped := s.data[0]
	*s = n

	return &popped
}

func (s *SearchQueue) Peek(i int) *SearchData {
	if i < len(s.data) {
		return &s.data[len(s.data)-i-1]
	}

	return nil
}

func (s *SearchQueue) Len() int {
	return len(s.data)
}