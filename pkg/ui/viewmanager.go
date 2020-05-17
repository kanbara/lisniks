package ui

import (
	"fmt"
	"github.com/awesome-gocui/gocui"
	"github.com/kanbara/lisniks/pkg/dictionary"
	"github.com/kanbara/lisniks/pkg/state"
	"github.com/sirupsen/logrus"
	"log"
)

type ViewManager struct {
	g     *gocui.Gui
	active []string
	currentView int

	Log   *logrus.Logger
	State *state.State
	Dict  *dictionary.Dictionary
	Views map[string]ViewUpdateSetter
}

func NewViewManager(dict *dictionary.Dictionary, state *state.State, log *logrus.Logger) *ViewManager {
	vm := &ViewManager{
		Log:   log,
		State: state,
		Dict:  dict,
	}

	// we can share this as a singleton as these functions are all nil
	// and all the updates are thread safe anyway
	dv := DefaultView{vm}

	vm.Views = map[string]ViewUpdateSetter{
		HeaderViewName: &HeaderView{dv},

		SearchViewName: &SearchView{View: View{
			ViewManager: vm,
			ViewsToUpdate: []string{
				POSViewName, LexViewName,
				CurrentWordViewName, LocalWordViewName,
				WordGrammarViewName, DefnViewName, StatusViewName}}},

		LexViewName: &LexiconView{ListView{
			View: View{
				ViewManager: vm,
				ViewsToUpdate: []string{
					POSViewName, CurrentWordViewName,
					LocalWordViewName, WordGrammarViewName,
					DefnViewName, StatusViewName}},
			viewName:     LexViewName,
			itemLen:      func() int { return len(vm.State.Words) },
			itemSelected: func() *int { return &vm.State.SelectedWord },
		}},

		CurrentWordViewName: &CurrentWordView{dv},
		LocalWordViewName:   &LocalWordView{dv},
		POSViewName:         &PartOfSpeechView{dv},
		WordGrammarViewName: &WordGrammarView{dv},
		DefnViewName:        &DefinitionView{dv},
		StatusViewName:      &StatusView{dv},
	}

	return vm
}

func (vm *ViewManager) Run() error {

	g, err := gocui.NewGui(gocui.Output256, false)
	if err != nil {
		log.Panicf("could not instantiate UI: %v", err)
	}

	vm.g = g
	defer vm.g.Close()

	vm.g.Highlight = true
	vm.g.SelFgColor = gocui.ColorGreen
	vm.g.SelFrameColor = gocui.ColorGreen

	for name, v := range vm.Views {
		if err := v.New(name); err != nil {
			vm.Log.Panicf("could instantiate view %v: %v", name, err)
		}

		if err := v.SetKeybindings(); err != nil {
			vm.Log.Panicf("could not set keybindings: %v", err)
		}
	}

	err = vm.SetGlobalKeybindings()
	if err != nil {
		vm.Log.Panicf("could not set keybinding: %v", err)
	}

	if err := vm.g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		// debug stuff if we crash sometimes
		// not sure if useful but oh well
		// TODO update this nonsense
		view := vm.g.CurrentView()
		ox, oy := view.Origin()
		vx, vy := view.Size()
		cx, cy := view.Cursor()
		cur, err := view.Line(cy)
		p := fmt.Sprintf("%v\nselected: %v\nview: %v\nview origin: %v,%v\n"+
			"view size: %v, %v\nview cursor: %v,%v\nlexicion list: %v\nbuf: `%v`",
			err, vm.State.SelectedWord, view.Name(), ox, oy, vx, vy, cx, cy, len(vm.State.Words), cur)
		panic(p)
	}

	return nil
}

func (vm *ViewManager) Layout(g *gocui.Gui) error {
	return nil
}

func (vm *ViewManager) SetActive(s ...string) {
	vm.active = s
}

func (vm *ViewManager) Cycle(updown int) error {
	// haha mod sucks https://github.com/golang/go/issues/448
	// ((m % n) + n) % n
	vm.currentView = ((vm.currentView + updown) + len(vm.active)) % len(vm.active)
	nextView := vm.active[vm.currentView]

	if _, err := vm.g.SetCurrentView(nextView); err != nil {
		return err
	}

	return nil
}

func (vm *ViewManager) NextView(g *gocui.Gui, _ *gocui.View) error {
	if err := vm.Cycle(1); err != nil {
		return err
	}

	return nil
}

func (vm *ViewManager) PrevView(g *gocui.Gui, _ *gocui.View) error {
	if err := vm.Cycle(-1); err != nil {
		return err
	}

	return nil
}

func (vm *ViewManager) SetGlobalKeybindings() error {
	if err := vm.g.SetKeybinding("",
		gocui.KeyCtrlC,
		gocui.ModNone,
		vm.QuitModal); err != nil {
		return err
	}

	if err := vm.g.SetKeybinding("",
		gocui.KeyCtrlR,
		gocui.ModNone,
		vm.ReloadModal); err != nil {
		return err
	}

	if err := vm.g.SetKeybinding("", '/', gocui.ModNone, vm.ToSearchView); err != nil {
		return err
	}

	if err := vm.g.SetKeybinding("", ']', gocui.ModNone, vm.NextView); err != nil {
		return err
	}

	if err := vm.g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, vm.NextView); err != nil {
		return err
	}

	if err := vm.g.SetKeybinding("", '[', gocui.ModNone, vm.PrevView); err != nil {
		return err
	}

	return nil
}

func (vm *ViewManager) UpdateStatusView() error {
	vm.g.Update(func(g *gocui.Gui) error {
		if v, err := g.View(StatusViewName); err != nil {
			return err
		} else {
			if err := vm.Views[StatusViewName].Update(v); err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

// make a viewPopped fn on the DefaultViews
func ToSearchView(g *gocui.Gui, _ *gocui.View) error {
	g.Cursor = true
	return ToView(g, SearchViewName)
}

func ToView(g *gocui.Gui, view string) error {
	_, err := g.SetCurrentView(view)
	if err != nil {
		return err
	}

	return nil
}
