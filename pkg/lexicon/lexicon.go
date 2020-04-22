package lexicon

import "github.com/kanbara/lisniks/pkg/word"

type Lexicon []word.Word

type Service struct {
	lexicon Lexicon
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

func NewLexiconService(lexicon Lexicon) *Service {
	return &Service{lexicon: lexicon}
}
