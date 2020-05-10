package dictionary

import (
	"fmt"
	"time"
)

func (d *Dictionary) Stats() string {
	return fmt.Sprintf("PolyGlot version %v, updated %v",
		d.file.Version, d.file.LastUpdated.Format(time.StampMilli))
}

func (d *Dictionary) LangAndVersion() string {
	return fmt.Sprintf("%v - %v",
		d.file.LanguageProperties.Name, d.file.LanguageProperties.Version())

}
