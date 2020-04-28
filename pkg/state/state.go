package state

import (
	"github.com/kanbara/lisniks/pkg/lexicon"
	"github.com/kanbara/lisniks/pkg/search"
	"github.com/kanbara/lisniks/pkg/word"
)

type State struct {
	Words        lexicon.Lexicon
	SelectedWord int
	StatusText   string
	HelpText     string

	SearchPattern  search.Pattern
	SearchPatterns map[search.Pattern]string
	SearchType     search.Type
	SearchTypes    map[search.Type]string
	SearchQueue    search.Queue
	CurrentSearch  search.Data
	QueuePos       int
}

func (s *State) CurrentWord() *word.Word {
	if len(s.Words) < s.SelectedWord {
		return nil
	}

	if len(s.Words) == 0 {
		return nil
	}

	return &s.Words[s.SelectedWord]
}
