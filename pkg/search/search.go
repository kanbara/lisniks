package search

import (
	"fmt"
	"strings"
)

type Type int
type Pattern int
type Identifier rune

const (
	TypeUnset Type = iota
	TypeAustrianWord
	TypeEnglishWord
	TypeWordDefinition
)

const (
	PatternUnset Pattern = iota
	PatternRegex
	PatternPhonotactic
)

const (
	IDAustrian Identifier = 'a'
	IDEnglish = 'e'
	IDDefinition = 'd'
	IDRegex = 'r'
	IDPhonotactic = 'p'
)

const RegexV = "[aeiouAEOIUāēīōūĀĒĪŌŪ]"
const RegexC = "[pbmvfhtdnszčǆšžňřłjrkgncPBMVFHTDNSZČǅǄŠŽŇŘŁJRKGNC]"

type Parsed struct {
	Type    Type
	Pattern Pattern
	String  string
}

type Queue struct {
	data []string
	max int
}

func TypeNames() map[Type]string {
	return map[Type]string{
		TypeEnglishWord: "english",
		TypeAustrianWord: "austrian",
		TypeWordDefinition: "definition",
	}
}

func PatternNames() map[Pattern]string {
	return map[Pattern]string{
		PatternRegex: "regex",
		PatternPhonotactic: "regex (phonotactic)",
	}
}

func parseTypeAndPattern(str string) (Type, Pattern, error) {
	t := TypeUnset
	p := PatternUnset
	// this is tricky, we need to check that we have the following patterns in the string
	// the only allowable bits are:
	// a e or d (one of type)
	// r or p (one of pattern)
	// both are optional, as the default should be `ar`
	// so the strings are: a,e,d,r,p, ar,er,dr, ap,ep,dp (and of course the other way around too)
	for _, c := range str {
		switch Identifier(c) {
		case IDAustrian:
			if t == TypeUnset {
				t = TypeAustrianWord
			} else {
				return TypeUnset, PatternUnset,
				fmt.Errorf("encountered %v but already set type to %v",
					string(c),
					TypeNames()[t])
			}
		case IDEnglish:
			if t == TypeUnset {
				t = TypeEnglishWord
			} else {
				return TypeUnset, PatternUnset,
					fmt.Errorf("encountered %v but already set type to %v",
						string(c),
						TypeNames()[t])
			}
		case IDDefinition:
			if t == TypeUnset {
				t = TypeWordDefinition
			} else {
				return TypeUnset, PatternUnset,
					fmt.Errorf("encountered %v but already set type to %v",
						string(c),
						TypeNames()[t])
			}
		case IDRegex:
			if p == PatternUnset {
				p = PatternRegex
			} else {
				return TypeUnset, PatternUnset,
					fmt.Errorf("encountered %v but already set type to %v",
						string(c),
						PatternNames()[p])


			}
		case IDPhonotactic:
			if p == PatternUnset {
				p = PatternPhonotactic
			} else {
				return TypeUnset, PatternUnset,
					fmt.Errorf("encountered %v but already set type to %v",
						string(c),
						PatternNames()[p])

			}
		default:
			return TypeUnset, PatternUnset, fmt.Errorf("encountered unknown symbol %v", string(c))
		}
	}

	if t == TypeUnset {
		t = TypeAustrianWord
	}

	if p == PatternUnset {
		p = PatternRegex
	}

	return t, p, nil
}

func ParseString(str string) (Parsed, error) {
	// oof this will be tricky -- the core idea is that we would have something like
	// TYPE_PATTERN/WORD/...
	// where i envision ... to be POS/SUBCAT (e.g. VERB/WEAK-O)
	// but we can also have the TYPE_PATTERN omitted, as the defaults (Austrian & Regex) are sensible
	// so we need to split and if we have two splits, for now, we assume the first is type and pattern
	// and the second is word
	// if we have one, we assume word
	//
	// well that's not so tricky but it can get more complicated in the future

	splits := strings.Split(str, "/")
	switch len(splits) {
	case 1: // just a regex, e.g. ^.*ǆa$
		return Parsed{
			Type:    TypeAustrianWord,
			Pattern: PatternRegex,
			String:  str,
		}, nil
	case 2: // a type+pattern and regex, e.g. e/^.*ǆa$
		t, p, err := parseTypeAndPattern(splits[0])
		return Parsed{
			Type:    t,
			Pattern: p,
			String:  splits[1],
		}, err
	}

	return Parsed{}, fmt.Errorf("search not parseable, check type, pattern, and /")
}

func NewQueue(max int) Queue {
	d := make([]string, 0, max)
	return Queue{max: max, data: d}
}

func (s *Queue) RemoveOthers(sd string) {
	n := 0
	nd := make([]string, len(s.data)) // at most we'll have len(String.data) points
	for _, x := range s.data {
		if x != sd {
			nd[n] = x // if x does not match sd, we keep it
			n++
		}
	}

	s.data = nd[:n]
}

func (s *Queue) Enqueue(sd string) {

	s.RemoveOthers(sd)

	if len(s.data)+1 > s.max {
		_ = s.Dequeue()
	}

	s.data = append(s.data, sd)
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