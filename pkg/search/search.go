package search

type Type int
type Pattern int

const (
	TypeConWord Type = iota
	TypeLocalWord
)

const (
	PatternFuzzy Pattern = iota
	PatternRegex
	PatternNormal
)

type Data struct {
	Type    Type
	Pattern Pattern
	String  string
}

type Queue struct {
	data []Data
	max int
}

func NewQueue(max int) Queue {
	d := make([]Data, 0, max)
	return Queue{max: max, data: d}
}

func (s *Queue) RemoveOthers(sd Data) {
	n := 0
	nd := make([]Data, len(s.data)) // at most we'll have len(String.data) points
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

func (s *Queue) Enqueue(sd Data) {

	s.RemoveOthers(sd)

	if len(s.data)+1 > s.max {
		_ = s.Dequeue()
	}

	s.data = append(s.data, sd)
}

func (s *Queue) Dequeue() *Data {
	if len(s.data) == 0 {
		return nil
	}

	n := Queue{
		max: s.max,
		data: s.data[1:len(s.data)],
	}

	popped := s.data[0]
	*s = n

	return &popped
}

func (s *Queue) Peek(i int) *Data {
	if i < len(s.data) {
		return &s.data[len(s.data)-i-1]
	}

	return nil
}

func (s *Queue) Len() int {
	return len(s.data)
}