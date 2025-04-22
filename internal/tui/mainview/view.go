package mainview

import (
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/glsubri/kex/pkg/dynlist"
)

type ItemStyles struct {
	Selected    lipgloss.Style
	Title       lipgloss.Style
	Type        lipgloss.Style
	Required    lipgloss.Style
	Description lipgloss.Style

	MdRenderer *MarkdownRenderer
}

func DefaultItemStyles() ItemStyles {
	styles := ItemStyles{
		Selected: lipgloss.NewStyle().
			BorderStyle(lipgloss.ThickBorder()).
			BorderLeft(true).
			BorderForeground(lipgloss.Color("62")),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Bold(true),
		Type: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		Required: lipgloss.NewStyle().
			Foreground(lipgloss.Color("167")).
			Bold(true),
		Description: lipgloss.NewStyle().
			Foreground(lipgloss.Color("248")),
	}

	styles.MdRenderer = NewMarkdownRenderer(80, styles.Description)

	return styles
}

type Styles struct {
	ItemStyles

	// Description and KGVK styles
	Title   lipgloss.Style
	Content lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		ItemStyles: DefaultItemStyles(),
		Title: lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1),
		Content: lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("250")),
	}
}

// Render implements the dynlist.Item interface for SchemaProperty
func (p SchemaProperty) Render(w io.Writer, state dynlist.ItemState) {

	// Render the title and type
	content := fmt.Sprintf("%s - %s",
		p.styles.Title.Render(p.Key),
		p.styles.Type.Render(fmt.Sprintf("%s", p.PropertyType)))

	// Add required indicator if needed
	if p.Required {
		content += " " + p.styles.Required.Render("[required]")
	}

	// Add description if any
	if p.Descr != "" {
		content += p.styles.MdRenderer.Render(p.Descr, state.Width)
	}

	// Apply selection style if selected
	if state.IsSelected {
		content = p.styles.Selected.PaddingLeft(1).Render(content)
	} else {
		content = lipgloss.NewStyle().PaddingLeft(2).Render(content)
	}

	fmt.Fprint(w, content)
}

// PrimaryAction implements the dynlist.Item interface for SchemaProperty
func (p SchemaProperty) PrimaryAction() tea.Cmd {
	if p.Ref == "" {
		return nil
	}

	return func() tea.Msg {
		return SelectPropertyMsg{Ref: p.Ref}
	}
}

func (m *Model) renderShortDescriptionSection() string {
	var content []string
	emptyLine := ""

	// Add KGVK section
	if kgvk := m.viewKGVK(); kgvk != "" {
		content = append(content, kgvk)
		content = append(content, emptyLine)
	}

	// Add Description section
	if desc := m.viewDescription(); desc != "" {
		content = append(content, desc)
		content = append(content, emptyLine)
	}

	// Limit the number of lines
	if len(content) > m.maxDescriptionLines {
		content = content[:m.maxDescriptionLines-1] // -1 to leave room for ellipsis
		content = append(content, "...")
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (m *Model) viewDescription() string {
	if m.schema.Description == "" {
		return ""
	}

	// Calculate available width for description
	dwx, _ := m.styles.Content.GetFrameSize()
	descriptionWidth := m.width - dwx

	// Wrap the description text
	wrapped := lipgloss.NewStyle().Width(descriptionWidth).Render(m.schema.Description)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.Title.Render("DESCRIPTION"),
		m.styles.Content.Render(wrapped),
	)
}

func (m *Model) viewKGVK() string {
	if len(m.schema.KGVKs) == 0 {
		return ""
	}

	const spacer = "   "

	viewLine := func(title, content string) string {
		return m.styles.Title.Render(title) + " " + m.styles.Content.Render(content) + spacer
	}
	singleView := func(k KGVK) string {
		var views []string

		if k.Group != "" {
			views = append(views, viewLine("GROUP", k.Group))
		}
		if k.Kind != "" {
			views = append(views, viewLine("KIND", k.Kind))
		}
		if k.Version != "" {
			views = append(views, viewLine("VERSION", k.Version))
		}

		return lipgloss.JoinHorizontal(lipgloss.Left, views...)
	}

	var view []string
	for _, kgvk := range m.schema.KGVKs {
		view = append(view, singleView(kgvk))
	}

	return lipgloss.JoinVertical(lipgloss.Top, view...)
}
