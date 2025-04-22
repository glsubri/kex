package mainview

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/glsubri/kex/pkg/dynlist"
)

// Schema represents a Kubernetes schema definition
type Schema struct {
	Fullpath    string
	Description string
	SchemaType  string
	Properties  []SchemaProperty
	KGVKs       []KGVK
}

// KGVK represents the metadata of top level kubernetes resources
type KGVK struct {
	Group   string
	Kind    string
	Version string
}

// SchemaProperty represents a property in a Kubernetes schema
type SchemaProperty struct {
	Key          string
	Descr        string
	PropertyType string
	Ref          string
	Required     bool

	styles ItemStyles
}

// Model is the main component for displaying schema information
type Model struct {
	width  int
	height int
	styles Styles

	schema Schema
	keys   keyMap

	// submodels
	propertiesList *dynlist.Model

	// description section settings
	maxDescriptionLines int
}

// New creates a new Model with the given schema and options
func New(schema Schema, opts ...Option) *Model {
	m := Model{
		schema:              schema,
		keys:                newKeyMap(),
		styles:              DefaultStyles(),
		maxDescriptionLines: 6,
	}

	// Apply options
	for _, opt := range opts {
		opt(&m)
	}

	// Create items for the list
	items := make([]dynlist.Item, len(schema.Properties))
	for i, prop := range schema.Properties {
		prop.styles = m.styles.ItemStyles
		items[i] = prop
	}

	m.propertiesList = dynlist.New(0, 0, items, dynlist.WithTitle("PROPERTIES"))

	return &m
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.up):
			m.propertiesList.MoveUp()
		case key.Matches(msg, m.keys.down):
			m.propertiesList.MoveDown()
		case key.Matches(msg, m.keys.enter):
			return m, m.propertiesList.ExecutePrimaryAction()
		}
	}

	return m, nil
}

// Help returns the keybindings for this model
func (m Model) Help() []key.Binding {
	return m.keys.ShortHelp()
}

// SetSize sets the size of the model
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	listHeight := height - lipgloss.Height(m.renderShortDescriptionSection())
	m.propertiesList.SetSize(width, listHeight)
}

// View returns the current view of the model
func (m *Model) View() string {
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderShortDescriptionSection(),
		m.propertiesList.View(),
	)

	return view
}

// SelectPropertyMsg is sent when a property with a ref is selected
type SelectPropertyMsg struct {
	Ref string
}
