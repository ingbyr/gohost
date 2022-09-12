package keys

import "github.com/charmbracelet/bubbles/key"

var (
	Up = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	)
	Down = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	)
	Left = key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "move left"),
	)
	Right = key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "move right"),
	)
	Help = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	)
	Quit = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)
	Esc = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "switch group or exit"),
	)
	Enter = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select or confirm"),
	)
	Switch = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch helpView"),
	)
	Save = key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	)
	New = key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "new"),
	)
)

func Arrows() []key.Binding {
	return []key.Binding{Up, Down, Left, Right}
}
