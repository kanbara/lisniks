package dictionary

type PartsOfSpeech []Part

type Part struct {
	ID int64 `xml:"classId"`
	Name string `xml:"className"`
	Notes Rawstring `xml:"classNotes"`
}

type POSMap map[int64]string

func MakePOSMap(pos PartsOfSpeech) POSMap {

	m := make(POSMap, len(pos))

	for _, p := range pos {
		m[p.ID] = p.Name
	}

	return m
}