package pathselection

import "strings"

type schemaItem struct {
	name     string
	fullpath string
}

func newItem(path string) schemaItem {

	return schemaItem{
		name:     schemaName(path),
		fullpath: path,
	}
}

func (s schemaItem) FilterValue() string {
	return s.fullpath
}

func (s schemaItem) Title() string {
	return s.name
}

func (s schemaItem) Description() string {
	return s.fullpath
}

func schemaName(path string) string {
	split := strings.Split(path, ".")
	return split[len(split)-1]
}
