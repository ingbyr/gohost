package tui

import "github.com/charmbracelet/bubbles/key"

var keys = newKeys()

type keyMaps struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Help   key.Binding
	Quit   key.Binding
	Esc    key.Binding
	Enter  key.Binding
	Switch key.Binding
	Save   key.Binding
	New    key.Binding
}

func newKeys() *keyMaps {
	return &keyMaps{
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
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "switch group or exit"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select or confirm"),
		),
		Switch: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "switch helpView"),
		),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		New: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "new"),
		),
	}
}

func (k keyMaps) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMaps) ArrowsHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Left, k.Right}
}

func (k keyMaps) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		//{k.Up, k.Down}, // column
		//{k.Left, k.Right},
		{k.Switch, k.New},
		{k.Save},
		{k.Help, k.Quit},
	}
}
