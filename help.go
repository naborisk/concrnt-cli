package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Unfocus    key.Binding
	Enter      key.Binding
	NewMessage key.Binding
	Refresh    key.Binding
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Help       key.Binding
	Quit       key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NewMessage, k.Refresh, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.NewMessage, k.Help, k.Quit},  // second column
	}
}

var keys = keyMap{
	Unfocus: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "unfocus"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "new message"),
	),
	NewMessage: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new message"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
