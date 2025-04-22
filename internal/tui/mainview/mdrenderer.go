package mainview

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const paddingAndBorder = 2

// MarkdownRenderer is used to render markdown. If no available markdown
// renderer is available, a fallback rendering is provided.
//
// Because of the way glamour renderer are created, we cannot change their
// width on the fly. MarkdownRenderer acts as a caching mechanism and
// recreates a new glamour renderer if needed.
type MarkdownRenderer struct {
	renderer       *glamour.TermRenderer
	fallbackStyle  lipgloss.Style
	width          int
	availableWidth int
}

func NewMarkdownRenderer(
	width int,
	fallbackStyle lipgloss.Style,
) *MarkdownRenderer {
	mdRenderer := MarkdownRenderer{
		renderer:       newGlamRenderer(width),
		fallbackStyle:  fallbackStyle,
		width:          width,
		availableWidth: width - 2,
	}

	return &mdRenderer
}

func newGlamRenderer(width int) *glamour.TermRenderer {
	glamourRenderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithEmoji(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil
	}

	return glamourRenderer
}

func (mdr *MarkdownRenderer) Render(in string, width int) string {
	// make sure that we have the correct width
	if mdr.width != width {
		mdr.width = width
		mdr.availableWidth = width - paddingAndBorder
		mdr.renderer = newGlamRenderer(mdr.availableWidth)
	}

	if mdr.renderer == nil {
		return mdr.renderFallback(in)
	} else {
		return mdr.renderGlamour(in)
	}
}

func (mdr *MarkdownRenderer) renderGlamour(in string) string {
	out, err := mdr.renderer.Render(in)
	if err != nil {
		return mdr.renderFallback(in)
	}

	// remove last '\n'
	out, _ = strings.CutSuffix(out, "\n")

	return out
}

func (mdr *MarkdownRenderer) renderFallback(in string) string {
	wrappedDesc := lipgloss.NewStyle().
		Width(mdr.availableWidth).
		Render(in)
	return "\n" + mdr.fallbackStyle.Render(wrappedDesc)
}
