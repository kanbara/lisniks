package dictionary

import (
	"fmt"
	"time"
)

func (d *Dictionary) Stats() string {
	return fmt.Sprintf("PolyGlot version %v, updated %v, word count %v",
		d.file.Version, d.file.LastUpdated.Format(time.StampMilli), len(d.file.Lexicon))
}

func (d *Dictionary) LangAndVersion() string {
	return fmt.Sprintf("%v - %v",
		d.file.LanguageProperties.Name, d.file.LanguageProperties.Version())

}
