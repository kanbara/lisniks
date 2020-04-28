package lexicon

import (
	"fmt"
	"github.com/kanbara/lisniks/pkg/language"
	"github.com/kanbara/lisniks/pkg/search"
	s "github.com/kanbara/lisniks/pkg/strings"
	"github.com/kanbara/lisniks/pkg/word"
	"regexp"
	"strings"
)

type Lexicon []word.Word

type Service struct {
	lexicon    Lexicon
	alphaOrder language.AlphaOrderMap
}

func (s *Service) GetByID(id int64) *word.Word {
	for _, i := range s.lexicon {
		if i.WordID == id {
			return &i
		}
	}

	return nil
}

func (s *Service) Words() []word.Word {
	return s.lexicon
}

func (s *Service) Count() int {
	return len(s.lexicon)
}

func (s *Service) At(index int) *word.Word {
	return &s.lexicon[index]
}

func (s *Service) found(str string, w s.Rawstring, pattern search.Pattern) (bool, error) {
	lstr := strings.ToLower(str)
	lw := strings.ToLower(string(w))

	switch pattern {
	case search.PatternFuzzy:
		if strings.Contains(lw, lstr) {
			return true, nil
		}
	case search.PatternNormal:
		if strings.HasPrefix(lw, lstr) {
			return true, nil
		}
	case search.PatternRegex:
		matched, err := regexp.Match(str, []byte(w))
		if err != nil {
			return false, err
		}

		return matched, nil
	}

	return false, nil
}

// TODO i have the idea that we should be able to chain search filters together
// like findByConWord(str).ByPartOfSpeech("verb").ByDefinitionContaining("foobar")
// which means we should have an output type that instead of returning []*word.Word
// should be another type like `Filtered` which is still just a []*word.Word
// and then all the searches are func (f *Filtered) ByFoo() *Filtered
// todo can also return the time or status string here to display to the view
func (s *Service) FindWords(str string, sp search.Pattern, st search.Type) (Lexicon, error) {
	// start with simple linear traversal here.
	// think about using suffix trees or something similar later,
	// or maybe rank queries with predecessor / successor
	// and fuzzy search with binary search
	var words []word.Word

	for i := range s.lexicon {
		switch st {
		case search.TypeConWord:
			if match, err := s.found(str, s.lexicon[i].Con, sp); err != nil {
				return nil, err
			} else if match {
				words = append(words, s.lexicon[i])
			}
		case search.TypeLocalWord:
			if match, err := s.found(str, s.lexicon[i].Local, sp); err != nil {
				return nil, err
			} else if match {
				words = append(words, s.lexicon[i])
			}
		}
	}

	return words, nil
}

func (s *Service) String() string {
	var out string
	for _, w := range s.lexicon {
		out += fmt.Sprintf("%v\n", w.Con)
	}

	return out
}

// TODO prolly should be a slice of rune instead
func NewLexiconService(lexicon Lexicon, alphaOrder language.AlphaOrderMap) *Service {
	return &Service{lexicon: lexicon, alphaOrder: alphaOrder}
}

// sort functions
// needed to make sort.Sort work on Service
// the only weird thing is less,
// because as opposed to normal strings, we need to sort on the Alpha Word Order
// from the dictionary, e.g. -AaĀāæ...Žž
// to ensure we order words the same as PolyGlot does

func (s Service) Len() int {
	return len(s.lexicon)
}

func (s Service) Swap(i, j int) {
	s.lexicon[i], s.lexicon[j] = s.lexicon[j], s.lexicon[i]
}

func (s Service) Less(i, j int) bool {

	// XXX stupid sorting, turns out this was the trick
	// PolyGlot strips spaces out of words when sorting
	runedI := []rune(strings.ReplaceAll(string(s.lexicon[i].Con), " ", ""))
	runedJ := []rune(strings.ReplaceAll(string(s.lexicon[j].Con), " ", ""))

	check := len(runedI)
	if len(runedI) > len(runedJ) {
		check = len(runedJ)
	}

	for i := 0; i < check; i++ {
		m := runedI[i]
		n := runedJ[i]

		diff := s.alphaOrder[m] - s.alphaOrder[n]
		switch {
		case diff < 0: // letter comes first
			return true
		case diff == 0: // letter is same
			continue
		case diff > 0: // letter comes after
			return false
		}
	}

	// at the very end, make sure we return the right one based on size
	// d.f. žoul comes before žoularis
	// and we checked all the letters and found them to match [ž, o, u l]
	return len(runedI) < len(runedJ)
}
