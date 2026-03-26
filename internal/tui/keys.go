package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Tab      key.Binding
	ShiftTab key.Binding
	Enter    key.Binding
	Add      key.Binding
	Remove   key.Binding
	Clear    key.Binding
	Inc      key.Binding
	Dec      key.Binding
	Search   key.Binding
	Reorder  key.Binding
	Back     key.Binding
	Quit     key.Binding
}

var keys = keyMap{
	Tab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
	ShiftTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev tab")),
	Enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Add:      key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add to cart")),
	Remove:   key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "remove")),
	Clear:    key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "clear cart")),
	Inc:      key.NewBinding(key.WithKeys("+"), key.WithHelp("+", "increase qty")),
	Dec:      key.NewBinding(key.WithKeys("-"), key.WithHelp("-", "decrease qty")),
	Search:   key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	Reorder:  key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "reorder")),
	Back:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}
