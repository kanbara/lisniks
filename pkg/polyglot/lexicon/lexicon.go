package lexicon

import (
	"fmt"
	"github.com/kanbara/lisniks/pkg/polyglot/language"
	"github.com/kanbara/lisniks/pkg/polyglot/word"
	"github.com/kanbara/lisniks/pkg/search"
	s "github.com/kanbara/lisniks/pkg/strings"
	"strings"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

type Lexicon []word.Word

type Service struct {
	lexicon    Lexicon
	alphaOrder language.AlphaOrderMap
}

func (se *Service) GetByID(id int64) *word.Word {
	for _, i := range se.lexicon {
		if i.WordID == id {
			return &i
		}
	}

	return nil
}

func (se *Service) Words() []word.Word {
	return se.lexicon
}

func (se *Service) Count() int {
	return len(se.lexicon)
}

func (se *Service) At(index int) *word.Word {
	return &se.lexicon[index]
}

func (se *Service) found(str string, w s.Rawstring, pattern search.Pattern) (bool, error) {
	switch pattern {
	case search.PatternRegex:
		re := pcre.MustCompile(str, pcre.MULTILINE|pcre.UTF8)
		matcher := re.Matcher([]byte(w.String()), 0)

		return matcher.Matches(), nil
	case search.PatternPhonotactic:
		// first substitute our C and V for the regex classes
		str = strings.ReplaceAll(str, "V", search.RegexV)
		str = strings.ReplaceAll(str, "C", search.RegexC)
		str := "^" + str + "$"

		re := pcre.MustCompile(str, pcre.MULTILINE|pcre.UTF8)
		matcher := re.Matcher([]byte(w.String()), 0)

		return matcher.Matches(), nil
	}

	return false, nil
}

// TODO i have the idea that we should be able to chain search filters together
// like findByConWord(str).ByPartOfSpeech("verb").ByDefinitionContaining("foobar")
// which means we should have an output type that instead of returning []*word.Word
// should be another type like `Filtered` which is still just a []*word.Word
// and then all the searches are func (f *Filtered) ByFoo() *Filtered
// todo can also return the time or status string here to display to the view
func (se *Service) FindWords(str string, posList map[int]bool) (Lexicon, error) {
	// start with simple linear traversal here.
	// think about using suffix trees or something similar later,
	// or maybe rank queries with predecessor / successor
	// and fuzzy search with binary search
	var words []word.Word

	parsed, err := search.ParseString(str)
	if err != nil {
		return words, err
	}

	var anyPOS int
	for _, v := range posList {
		if v {
			anyPOS++
		}
	}

	for i := range se.lexicon {
		var r s.Rawstring
		switch parsed.Type {
		case search.TypeAustrianWord:
			r = se.lexicon[i].Austrian
		case search.TypeEnglishWord:
			r = se.lexicon[i].English
		case search.TypeWordDefinition:
			r = se.lexicon[i].Def
		}

		match, err := se.found(parsed.String, r, parsed.Pattern)
		if err != nil {
			return nil, err
		}

		// gotta check if anyPOS of the POSes are set
		// if none are set, we check filter, otherwise we add the word


		if match {
			// we have no filters, add the word
			if anyPOS == 0 {
				words = append(words, se.lexicon[i])
				continue
			}

			if isSet, ok := posList[int(se.lexicon[i].Type)]; ok && isSet {
				words = append(words, se.lexicon[i])
			}
		}
	}

	return words, nil
}

func (se *Service) String() string {
	var out string
	for _, w := range se.lexicon {
		out += fmt.Sprintf("%v\n", w.Austrian)
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

func (se Service) Len() int {
	return len(se.lexicon)
}

func (se Service) Swap(i, j int) {
	se.lexicon[i], se.lexicon[j] = se.lexicon[j], se.lexicon[i]
}

func (se Service) Less(i, j int) bool {

	// XXX stupid sorting, turns out this was the trick
	// PolyGlot strips spaces out of words when sorting
	runedI := []rune(strings.ReplaceAll(string(se.lexicon[i].Austrian), " ", ""))
	runedJ := []rune(strings.ReplaceAll(string(se.lexicon[j].Austrian), " ", ""))

	check := len(runedI)
	if len(runedI) > len(runedJ) {
		check = len(runedJ)
	}

	for i := 0; i < check; i++ {
		m := runedI[i]
		n := runedJ[i]

		diff := se.alphaOrder[m] - se.alphaOrder[n]
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
