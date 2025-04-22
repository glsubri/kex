package pathselection

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

type Model struct {
	width  int
	height int

	keys         keyMap
	internalList list.Model
}

func New(schemas []string) *Model {
	items := make([]list.Item, len(schemas))
	for i, s := range schemas {
		items[i] = newItem(s)
	}

	l := list.New(items, list.NewDefaultDelegate(), 40, 40)
	l.SetShowHelp(false)
	l.Title = "Schemas"
	l.Filter = customFilter
	l.SetFilterState(list.Filtering)

	m := Model{
		internalList: l,
		keys:         newKeyMap(),
	}

	return &m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle top-level model messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't interfere when filtering
		if m.internalList.FilterState() == list.Filtering {
			break
		}

		switch {
		// Item selection
		case key.Matches(msg, m.keys.enter):
			item, ok := m.internalList.SelectedItem().(schemaItem)
			if !ok {
				panic("selected item was not of expected type 'schemaItem'")
			}

			return m, func() tea.Msg { return SelectedItemMsg{item.fullpath} }
		// Quit view
		case key.Matches(msg, m.keys.quit):
			return m, func() tea.Msg { return QuitViewMsg{} }
		}
	}

	// Propagate update to internal list by default
	var cmd tea.Cmd
	m.internalList, cmd = m.internalList.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Left,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.internalList.View(),
		),
	)
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.internalList.SetSize(width, height)
}

func (m *Model) OpenToFiltering() {
	m.internalList.SetFilterState(list.Filtering)
}

// Help returns the keybindings for this model
func (m *Model) Help() []key.Binding {
	return m.keys.ShortHelp()
}

// Msg & Cmd

type SelectedItemMsg struct {
	Path string
}

type QuitViewMsg struct{}

// Internal

func customFilter(term string, targets []string) []list.Rank {
	matches := fuzzy.FindNoSort(term, targets)
	sort.SliceStable(matches, func(i, j int) bool {
		wi := targets[matches[i].Index]
		schemaNameI := schemaName(wi)
		wj := targets[matches[j].Index]
		schemaNameJ := schemaName(wj)

		// Exact match on fullpath
		if wi == term && wj != term {
			return true
		}
		if wj == term && wi != term {
			return false
		}

		// Exact match on schema name
		if schemaNameI == term && schemaNameJ != term {
			return true
		}
		if schemaNameJ == term && schemaNameI != term {
			return false
		}

		// Schema name has prefix
		if strings.HasPrefix(schemaNameI, term) && !strings.HasPrefix(schemaNameJ, term) {
			return true
		}
		if strings.HasPrefix(schemaNameJ, term) && !strings.HasPrefix(schemaNameI, term) {
			return false
		}

		// Fallback to size of string to target simpler structs
		//
		// Filtering for "Replica", we should prioritize "ReplicaSet" or "ReplicaSetSpec"
		// over "ReplicationControllerStatus".
		//
		// It will not lessen the UX because the user can simply type "ReplicaContrStat" or
		// something similair to select only the end of the word, thanks to fuzzy finding.
		return len(schemaNameI) <= len(schemaNameJ)
	})

	result := make([]list.Rank, len(matches))
	for i, m := range matches {
		result[i] = list.Rank{
			Index:          m.Index,
			MatchedIndexes: m.MatchedIndexes,
		}
	}
	return result
}
