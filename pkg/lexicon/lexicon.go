package lexicon

import (
	"fmt"
	"github.com/kanbara/lisniks/pkg/language"
	"github.com/kanbara/lisniks/pkg/word"
	"strings"
)

type Lexicon []word.Word

type Service struct {
	lexicon Lexicon
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

func (s *Service) Count() int {
	return len(s.lexicon)
}

func (s *Service) At(index int) *word.Word {
	return &s.lexicon[index]
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