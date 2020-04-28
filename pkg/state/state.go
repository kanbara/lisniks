package state

import (
	"github.com/kanbara/lisniks/pkg/lexicon"
	"github.com/kanbara/lisniks/pkg/word"
)

type State struct {
	Words        lexicon.Lexicon
	SelectedWord int
	StatusText   string
	HelpText     string

	SearchPattern  lexicon.SearchPattern
	SearchPatterns map[lexicon.SearchPattern]string
	SearchType     lexicon.SearchType
	SearchTypes    map[lexicon.SearchType]string
	SearchQueue     SearchQueue
	CurrentSearch string
	QueuePos int
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
