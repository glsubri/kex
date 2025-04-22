package tui

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/glsubri/kex/internal/tui/loadingview"
	"github.com/glsubri/kex/internal/tui/mainview"
	"github.com/glsubri/kex/internal/tui/pathselection"
	"github.com/go-openapi/spec"
)

const (
	kubernetesKGVKKey = "x-kubernetes-group-version-kind"
)

type pageState int

const (
	loadingView pageState = iota
	pathSelectionView
	mainView
)

type Model struct {
	// UI State
	width       int
	height      int
	currentPage pageState

	// Submodels
	loadingModel       *loadingview.Model
	pathSelectionModel *pathselection.Model
	mainViewModel      *mainview.Model

	// Main model
	currentPath string
	spec        *spec.Swagger
	keys        keyMap
	history     *History
}

func NewModel(swagger *spec.Swagger) *Model {
	m := Model{
		width:       0,
		height:      0,
		currentPage: loadingView,

		loadingModel:       loadingview.New(),
		pathSelectionModel: nil,
		mainViewModel:      nil,

		currentPath: "",
		spec:        swagger,
		keys:        newKeyMap(),
		history:     NewHistory(),
	}

	return &m
}

type initMsg struct {
	err  error
	spec *spec.Swagger
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		//submodels
		m.loadingModel.Init(),
		// commands
		loadOpenAPISpecFromKubectl(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Submodel signals and top level controls
	switch msg := msg.(type) {

	// Initial msg
	case initMsg:
		if msg.err != nil {
			return m, tea.Quit
		}

		m.spec = msg.spec

		schemas := make([]string, 0)
		for path := range m.spec.Definitions {
			schemas = append(schemas, path)
		}

		m.currentPage = pathSelectionView
		m.pathSelectionModel = pathselection.New(schemas)
		m.updateSubmodelSizes()
		return m, nil

	// Path selection Msg
	case pathselection.SelectedItemMsg:
		if !m.changeToSchema(msg.Path) {
			return m, nil
		}
		m.history.Clear()
		m.history.Add(msg.Path)
		return m, nil

	case pathselection.QuitViewMsg:
		// consider as noop if no previous data
		if m.currentPath == "" {
			return m, nil
		}

		m.currentPage = mainView
		m.history.Clear()
		return m, nil

	// Description Msg
	case mainview.SelectPropertyMsg:
		if !m.changeToSchema(msg.Ref) {
			return m, nil
		}
		m.history.Add(msg.Ref)
		return m, nil

	// Global keys
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.open):
			if m.currentPage == mainView {
				m.changePageToSelectionView()
				return m, nil
			}
		case key.Matches(msg, m.keys.back):
			if m.currentPage == mainView && len(m.history.items) > 1 {
				m.history.Pop()
				previousRef := m.history.items[len(m.history.items)-1].Ref

				m.changeToSchema(previousRef)
				return m, nil
			}
		}

	// Size msg
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSubmodelSizes()
		return m, nil
	}

	// Submodels messages
	switch m.currentPage {
	case loadingView:
		newTeaModel, newCmd := m.loadingModel.Update(msg)
		newModel, ok := newTeaModel.(*loadingview.Model)
		if !ok {
			panic("could not cast loadingview.model")
		}
		m.loadingModel = newModel
		return m, newCmd

	case pathSelectionView:
		newTeaModel, newCmd := m.pathSelectionModel.Update(msg)
		newModel, ok := newTeaModel.(*pathselection.Model)
		if !ok {
			panic("could not cast pathselection.model")
		}
		m.pathSelectionModel = newModel
		return m, newCmd

	case mainView:
		newTeaModel, newCmd := m.mainViewModel.Update(msg)
		newModel, ok := newTeaModel.(*mainview.Model)
		if !ok {
			panic("could not cast mainview.model")
		}
		m.mainViewModel = newModel
		return m, newCmd
	}

	return m, nil
}

// internal methods

func (m *Model) globalHelp() []key.Binding {
	keys := []key.Binding{m.keys.quit}

	if m.currentPage == mainView {
		keys = append(
			keys,
			m.keys.open,
			m.keys.back,
		)
	}

	return keys
}

func (m *Model) changePageToSelectionView() {
	m.pathSelectionModel.OpenToFiltering()
	m.currentPage = pathSelectionView
}

func (m *Model) updateSubmodelSizes() {
	aboveHeight := lipgloss.Height(m.titleView())
	belowHeight := lipgloss.Height(m.helpView())

	px, py := panelBox.GetFrameSize()
	contentWidth := m.width - px
	contentHeight := m.height - py - aboveHeight - belowHeight - 1

	if m.loadingModel != nil {
		m.loadingModel.SetSize(contentWidth, contentHeight)
	}
	if m.pathSelectionModel != nil {
		m.pathSelectionModel.SetSize(contentWidth, contentHeight)
	}
	if m.mainViewModel != nil {
		m.mainViewModel.SetSize(contentWidth, contentHeight)
	}
}

func (m *Model) changeToSchema(ref string) bool {
	schema, exists := m.spec.Definitions[ref]
	if !exists {
		return false
	}

	internalSchema := fromSpecSchema(ref, &schema)
	m.mainViewModel = mainview.New(internalSchema)
	m.currentPath = ref
	m.currentPage = mainView
	m.updateSubmodelSizes()
	return true
}

// internal functions

func loadOpenAPISpecFromKubectl() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("kubectl", "get", "--raw", "/openapi/v2")
		output, err := cmd.Output()
		if err != nil {
			return initMsg{err: err}
		}

		var swagger spec.Swagger
		if err := json.Unmarshal(output, &swagger); err != nil {
			return initMsg{err: err}
		}

		return initMsg{spec: &swagger}
	}
}

func fromSpecSchema(fullpath string, s *spec.Schema) mainview.Schema {
	return mainview.Schema{
		Fullpath:    fullpath,
		Description: s.Description,
		SchemaType:  strings.Join(s.Type, ", "),
		Properties:  extractProperties(s),
		KGVKs:       extractKGVKs(s),
	}
}

func extractKGVKs(s *spec.Schema) []mainview.KGVK {
	ext, ok := s.Extensions[kubernetesKGVKKey]
	if !ok {
		return nil
	}

	raw, err := json.Marshal(ext)
	if err != nil {
		return nil
	}

	var result []mainview.KGVK
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil
	}

	return result
}

func extractProperties(s *spec.Schema) []mainview.SchemaProperty {
	props := make([]mainview.SchemaProperty, 0, len(s.Properties))

	requiredSet := map[string]struct{}{}
	for _, r := range s.Required {
		requiredSet[r] = struct{}{}
	}

	for key, prop := range s.Properties {
		_, isRequired := requiredSet[key]

		// Start with basic type description
		var (
			propertyType string
			ref          string
			err          error
		)

		switch {
		// silently an object
		case len(prop.Type) == 0:
			propertyType, ref, err = handleObjects(prop)
		case prop.Type[0] == "array":
			propertyType, ref, err = handleArrays(prop)
		case prop.Type[0] == "object":
			propertyType, ref, err = handleObjects(prop)
		default:
			propertyType = prop.Type[0]
		}
		if err != nil {
			panic(err)
		}

		props = append(props, mainview.SchemaProperty{
			Key:          key,
			Descr:        prop.Description,
			PropertyType: propertyType,
			Ref:          ref,
			Required:     isRequired,
		})
	}

	// Sort properties by key
	sort.Slice(props, func(i, j int) bool {
		return props[i].Key < props[j].Key
	})

	return props
}

func cleanRef(ref string) string {
	return strings.TrimPrefix(ref, "#/definitions/")
}

func displayableRef(ref string) string {
	cleanRef := cleanRef(ref)
	split := strings.Split(cleanRef, ".")

	return split[len(split)-1]
}

func handleArrays(prop spec.Schema) (string, string, error) {
	if prop.Items.Schema == nil {
		return "", "", fmt.Errorf("no schema found for array")
	}

	// If the array items have a reference, handle that
	if ref := prop.Items.Schema.Ref.String(); ref != "" {
		cleanedRef := cleanRef(ref)
		displayName := displayableRef(ref)
		return fmt.Sprintf("[]%s", displayName), cleanedRef, nil
	}

	// If it's an array of primitive types
	if len(prop.Items.Schema.Type) > 0 {
		return fmt.Sprintf("[]%s", strings.Join(prop.Items.Schema.Type, ", ")), "", nil
	}

	return "", "", fmt.Errorf("invalid array schema")
}

func handleObjects(prop spec.Schema) (string, string, error) {
	ref := prop.Ref.String()

	// If it's an object type without a reference
	if ref == "" {
		return "object", "", nil
	}

	// If the object has a reference, use that
	cleanedRef := cleanRef(ref)
	displayName := displayableRef(ref)
	return displayName, cleanedRef, nil
}
