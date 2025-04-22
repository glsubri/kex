package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	quit key.Binding
	open key.Binding
	back key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit program"),
		),
		open: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open selection menu"),
		),
		back: key.NewBinding(
			key.WithKeys("ctrl+o"),
			key.WithHelp("ctrl+o", "go back"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit, k.open, k.back}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.quit, k.open, k.back}}
}
