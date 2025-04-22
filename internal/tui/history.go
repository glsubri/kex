package tui

import "strings"

type historyItem struct {
	Ref      string
	ShortRef string
}

type History struct {
	items []historyItem
}

func NewHistory() *History {
	return &History{
		items: make([]historyItem, 0),
	}
}

func (h *History) Add(ref string) {
	shortRef := getShortRef(ref)
	h.items = append(h.items, historyItem{Ref: ref, ShortRef: shortRef})
}

func (h *History) Pop() {
	if len(h.items) > 0 {
		h.items = h.items[:len(h.items)-1]
	}
}

func (h *History) Clear() {
	h.items = make([]historyItem, 0)
}

func getShortRef(ref string) string {
	parts := strings.Split(ref, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}
