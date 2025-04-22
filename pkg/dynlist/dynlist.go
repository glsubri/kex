// dynlist is a bubbles list component that can have elements of different heights
package dynlist

import (
	"fmt"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Title lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1),
	}
}

const (
	sameItemScrollAmount = 4
)

type ItemState struct {
	// Available width for the item.
	Width int
	// Whether the item is selected.
	IsSelected bool
}

// Item is an item in the list.
type Item interface {
	// Render renders the item.
	Render(w io.Writer, state ItemState)
	// PrimaryAction returns a command to execute
	PrimaryAction() tea.Cmd
}

type Model struct {
	Title     string
	ShowTitle bool

	Styles Styles

	width  int
	height int

	items        []Item
	selectedItem int

	visible section

	// cache section
	cache          []string
	cacheIsInvalid bool
	itemsSection   map[int]section
}

func New(
	width int,
	height int,
	items []Item,
	opts ...Option,
) *Model {
	m := Model{
		Title:     "",
		ShowTitle: false,

		Styles: DefaultStyles(),

		width:  width,
		height: height,

		items:        items,
		selectedItem: 0,

		visible: section{start: 0, end: height},

		// cache section
		cache:          make([]string, 0),
		cacheIsInvalid: true,
		itemsSection:   make(map[int]section),
	}

	// Apply options
	for _, opt := range opts {
		opt(&m)
	}

	return &m
}

// View returns the current view of the list component.
func (m *Model) View() string {
	// Don't even render if there are no items
	if len(m.items) == 0 {
		return ""
	}

	// Populate cache if needed
	if m.cacheIsInvalid {
		m.populateCache()
	}

	// Build the view
	var s strings.Builder
	if m.ShowTitle {
		s.WriteString(m.renderTitle())
	}

	start := m.visible.start
	end := min(m.visible.end, len(m.cache))
	for i := start; i < end; i++ {
		s.WriteString(m.cache[i])
		s.WriteString("\n")
	}

	return s.String()
}

// SetTitle sets the title of the list.
func (m *Model) SetTitle(title string) {
	m.Title = title
}

// SetShowTitle sets whether the title should be shown.
func (m *Model) SetShowTitle(show bool) {
	m.ShowTitle = show
}

// SetSize sets the width and height of this component.
func (m *Model) SetSize(width, height int) {
	m.height = height
	m.width = width

	availableHeight := height - lipgloss.Height(m.renderTitle())

	m.visible = section{
		start: m.visible.start,
		end:   m.visible.start + availableHeight - 1, // -1 because we don't want to include the last "\n"
	}
	m.cacheIsInvalid = true
}

// MoveDown moves the selection down. If necessary, it will scroll the list.
func (m *Model) MoveDown() {
	if len(m.items) == 0 {
		return
	}

	// if the current item is not fully shown, don't change item, just scroll
	itemAIndex := m.selectedItem
	itemA := m.itemsSection[itemAIndex]
	if itemA.end > m.visible.end {
		m.viewportScrollDown(sameItemScrollAmount)
		return
	}

	// otherwise change item
	m.selectedItem = clamp(m.selectedItem+1, 0, len(m.items)-1)
	itemBIndex := m.selectedItem

	// refresh UI (for selection state)
	m.cacheItem(itemAIndex)
	m.cacheItem(itemBIndex)

	itemB := m.itemsSection[itemBIndex]
	// we don't need to scroll if above center line
	if itemB.center() < m.visible.center() {
		return
	}

	// gentle scroll if cannot be made fully visible
	if itemB.height() > m.visible.height() {
		m.viewportScrollDown(sameItemScrollAmount)
		return
	}

	// item can be shown fully, so align to center
	scrollAmount := itemB.center() - m.visible.center() + 1
	m.viewportScrollDown(scrollAmount)
}

// MoveUp moves the selection up. If necessary, it will scroll the list.
func (m *Model) MoveUp() {
	if len(m.items) == 0 {
		return
	}

	// if the current item is not fully shown, don't change item, just scroll
	itemAIndex := m.selectedItem
	itemA := m.itemsSection[itemAIndex]
	if itemA.start < m.visible.start {
		m.viewportScrollUp(sameItemScrollAmount)
		return
	}

	// otherwise change item
	m.selectedItem = clamp(m.selectedItem-1, 0, len(m.items)-1)
	itemBIndex := m.selectedItem

	// refresh UI (for selection state)
	m.cacheItem(itemAIndex)
	m.cacheItem(itemBIndex)

	itemB := m.itemsSection[itemBIndex]
	// we don't need to scroll if below center line
	if itemB.center() > m.visible.center() {
		return
	}

	// gentle scroll if cannot be made fully visible
	if itemB.height() > m.visible.height() {
		m.viewportScrollUp(sameItemScrollAmount)
		return
	}

	// item can be shown fully, so align to center
	scrollAmount := m.visible.center() - itemB.center() + 1
	m.viewportScrollUp(scrollAmount)
}

func (m *Model) ExecutePrimaryAction() tea.Cmd {
	item := m.items[m.selectedItem]
	return item.PrimaryAction()
}

// internal

func (m *Model) populateCache() {
	var cache []string
	var currentLine int

	// Render each item and split by newlines
	for i := range m.items {
		lines := m.renderItem(i)

		// Update cache
		cache = append(cache, lines...)

		// Update item's section metadata
		start := currentLine
		end := len(lines) + currentLine
		m.itemsSection[i] = section{start, end}

		currentLine = end
	}

	m.cache = cache
	m.cacheIsInvalid = false
}

func (m *Model) viewportScrollUp(n int) {
	m.visible = m.visible.scrollUp(n)
}

func (m *Model) viewportScrollDown(n int) {
	m.visible = m.visible.scrollDown(n, len(m.cache))
}

func (m *Model) renderTitle() string {
	return m.Styles.Title.Render(m.Title) + "\n\n"
}

func (m *Model) renderItem(index int) []string {
	item := m.items[index]
	state := ItemState{
		Width:      m.width,
		IsSelected: index == m.selectedItem,
	}

	var s strings.Builder
	item.Render(&s, state)
	s.WriteString("\n")

	// Split the rendered content by newlines and add each line to cache
	str := s.String()
	str = strings.Replace(str, "\n\r", "\n", -1)
	lines := strings.Split(str, "\n")

	return lines
}

func (m *Model) cacheItem(index int) {
	originalSection := m.itemsSection[index]
	data := m.renderItem(index)

	if len(data) != originalSection.height() {
		panic(fmt.Sprintf(
			"caching item %d: new data has height %d, need %d\ndata:\n%+v\nsection:\n%+v",
			index,
			len(data),
			originalSection.height(),
			data,
			originalSection,
		))
	}

	for i, line := range data {
		cacheIndex := i + originalSection.start
		m.cache[cacheIndex] = line
	}
}

func clamp(v, low, high int) int {
	if low > high {
		low, high = high, low
	}
	return max(low, min(v, high))
}
