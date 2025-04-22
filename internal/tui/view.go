package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Align(lipgloss.Center)

	panelBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
		// Margin(1, 2).
		Padding(1, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Align(lipgloss.Center)

	historyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 2)
)

func (m *Model) View() string {
	var content string

	switch m.currentPage {
	case pathSelectionView:
		content = panelBox.Render(m.pathSelectionModel.View())
	case mainView:
		content = panelBox.Render(m.mainViewModel.View())
	case loadingView:
		content = m.loadingModel.View()
	default:
		content = "incomplete switch: no selected view"
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.titleView(),
		m.historyView(),
		content,
		m.helpView(),
	)
}

func (m *Model) titleView() string {
	return titleStyle.Render("kex: 📦 ", m.currentPath)
}

func (m *Model) helpView() string {
	// global help
	help := m.globalHelp()

	// Add current view's keys
	switch m.currentPage {
	case pathSelectionView:
		help = append(help, m.pathSelectionModel.Help()...)
	case mainView:
		if m.mainViewModel == nil {
			panic("invariant violation: mainViewModel must not be nil when in mainView")
		}
		help = append(help, m.mainViewModel.Help()...)
	}

	// Format help text
	helpText := make([]string, 0, len(help))
	for _, binding := range help {
		if binding.Help().Key != "" {
			helpText = append(helpText, binding.Help().Key+": "+binding.Help().Desc)
		}
	}

	return helpStyle.Width(m.width).Render(strings.Join(helpText, " • "))
}

func (m *Model) historyView() string {
	l := len(m.history.items)
	switch {
	case l == 0:
		return ""
	case l <= 4:
		items := make([]string, len(m.history.items))
		for i, item := range m.history.items {
			items[i] = item.ShortRef
		}
		return historyStyle.Render(strings.Join(items, " > "))
	default:
		line := fmt.Sprintf("[%d] > %s > %s", l-2, m.history.items[l-2].ShortRef, m.history.items[l-1].ShortRef)
		return historyStyle.Render(line)
	}
}
