package pathselection

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	quit  key.Binding
	enter key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit selection view"),
		),
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit, k.enter}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.quit, k.enter}}
}
