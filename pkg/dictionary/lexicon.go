package dictionary

type Lexicon []Word

func (l Lexicon) GetByID(id int64) *Word {
	for _, i := range l {
		if i.WordID == id {
			return &i
		}
	}

	return nil
}