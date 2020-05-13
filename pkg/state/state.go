package state

import (
	"github.com/kanbara/lisniks/pkg/dictionary"
	"github.com/kanbara/lisniks/pkg/lexicon"
	"github.com/kanbara/lisniks/pkg/search"
	"github.com/kanbara/lisniks/pkg/word"
)

type State struct {
	Version      string
	Words        lexicon.Lexicon
	SelectedWord int
	StatusText   string
	HelpText     string

	SearchPattern  search.Pattern
	SearchPatterns map[search.Pattern]string
	SearchType     search.Type
	SearchTypes    map[search.Type]string
	SearchQueue    search.Queue
	CurrentSearch  string
	QueuePos       int
}

func NewState(version string, dict *dictionary.Dictionary) *State {
	return &State{
		Version: version,
		Words:        dict.Lexicon.Words(),
		SelectedWord: 0,
		SearchPattern: search.PatternRegex,
		SearchPatterns: map[search.Pattern]string{
			search.PatternRegex:       search.PatternNames()[search.PatternRegex],
			search.PatternPhonotactic: search.PatternNames()[search.PatternPhonotactic],
		},
		SearchType: search.TypeAustrianWord,
		SearchTypes: map[search.Type]string{
			search.TypeAustrianWord:   search.TypeNames()[search.TypeAustrianWord],
			search.TypeEnglishWord:    search.TypeNames()[search.TypeEnglishWord],
			search.TypeWordDefinition: search.TypeNames()[search.TypeWordDefinition],
		},
		SearchQueue: search.NewQueue(50),
		QueuePos: -1,
	}
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
