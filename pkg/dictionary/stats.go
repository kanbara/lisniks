package dictionary

import "fmt"

func (d *Dictionary) Stats() string {
	return fmt.Sprintf("loaded dictionary from PolyGlot version %v, updated %v, word count %v",
		d.file.Version, d.file.LastUpdated, len(d.file.Lexicon))
}

func (d *Dictionary) LangAndVersion() string {
	return fmt.Sprintf("%v - %v\n\n",
		d.file.LanguageProperties.Name, d.file.LanguageProperties.Version())

}
