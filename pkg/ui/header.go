package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
)

func (m *Manager) NewHeaderView(g* gocui.Gui) error {
	stats := m.dict.Stats()
	langAndVersion := m.dict.LangAndVersion()
	stringlen := len(stats)

	// TODO move positions where possible into views.go maybe
	if v, err := g.SetView(langView, 0, 0, stringlen+1, 5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		_, err := fmt.Fprintln(v, langAndVersion)
		_, err = fmt.Fprintln(v, stats)

		if err != nil {
			return err
		}

		v.Frame = false
	}

	return nil
}
