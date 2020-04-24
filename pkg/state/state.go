package state

import (
	"github.com/kanbara/lisniks/pkg/lexicon"
	"github.com/kanbara/lisniks/pkg/word"
)

type State struct {
	Words lexicon.Lexicon
	SelectedWord int
}

func (s *State) CurrentWord() *word.Word {
	return &s.Words[s.SelectedWord]
}