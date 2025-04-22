package mainview

import "github.com/charmbracelet/lipgloss"

// Option is a function that configures a Model
type Option func(*Model)

// WithStyles sets custom styles for the model
func WithStyles(styles Styles) Option {
	return func(m *Model) {
		m.styles = styles
	}
}

// WithItemStyles sets custom item styles for the model
func WithItemStyles(styles ItemStyles) Option {
	return func(m *Model) {
		m.styles.ItemStyles = styles
	}
}

// WithTitleStyle sets a custom title style
func WithTitleStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.styles.Title = style
	}
}

// WithContentStyle sets a custom content style
func WithContentStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.styles.Content = style
	}
}

// WithSelectedStyle sets a custom selected style for items
func WithSelectedStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.styles.ItemStyles.Selected = style
	}
}

// WithItemTitleStyle sets a custom title style for items
func WithItemTitleStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.styles.ItemStyles.Title = style
	}
}

// WithItemTypeStyle sets a custom type style for items
func WithItemTypeStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.styles.ItemStyles.Type = style
	}
}

// WithItemRequiredStyle sets a custom required style for items
func WithItemRequiredStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.styles.ItemStyles.Required = style
	}
}

// WithItemDescriptionStyle sets a custom description style for items
func WithItemDescriptionStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.styles.ItemStyles.Description = style
	}
}

// WithMaxDescriptionLines sets the maximum number of lines to show in the description section
func WithMaxDescriptionLines(lines int) Option {
	return func(m *Model) {
		m.maxDescriptionLines = lines
	}
}
