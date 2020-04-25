package ui

const (
	langView = "lang"
	lexView = "lexicon"
	posView = "part of speech"
	wordGrammarView = "word classes"
	localWordView = "local word"
	defnView = "definition"
)

//func (m *Manager) NextView(g *gocui.Gui, v *gocui.View) error {
//	a.viewIndex = (a.viewIndex + 1) % len(VIEWS)
//	return a.setView(g)
//}
//
//func (m *Manager) PrevView(g *gocui.Gui, v *gocui.View) error {
//	a.viewIndex = (a.viewIndex - 1 + len(VIEWS)) % len(VIEWS)
//	return a.setView(g)
//}

//func (m *Manager) setView(g *gocui.Gui) error {
//	_, err := g.SetCurrentView(lexView)
//	return err
//}