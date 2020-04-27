package state

type data []string
type Queue struct {
	data
	max int
}

func NewQueue(max int) Queue {
	d := make(data, 0, max)
	return Queue{max: max, data: d}
}

func (s *Queue) RemoveOthers(str string) {
	n := 0
	for _, x := range s.data {
		if x != str {
			s.data[n] = x
			n++
		}
	}
	s.data = s.data[:n]
}

func (s *Queue) Enqueue(str string) {

	s.RemoveOthers(str)
	if len(s.data)+1 > s.max {
		_ = s.Dequeue()
	}

	s.data = append(s.data, str)
}

func (s *Queue) Dequeue() *string {
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

func (s *Queue) Peek(i int) *string {
	if i < len(s.data) {
		return &s.data[len(s.data)-i-1]
	}

	return nil
}

func (s *Queue) Len() int {
	return len(s.data)
}