package dynlist

// Option is a function that configures a Model.
type Option func(*Model)

// WithTitle sets the title of the list.
// By default, if the title is set, it will be shown.
func WithTitle(title string) Option {
	return func(m *Model) {
		m.Title = title
		m.ShowTitle = true
	}
}

// WithStyles sets the styles of the list.
func WithStyles(styles Styles) Option {
	return func(m *Model) {
		m.Styles = styles
	}
}
