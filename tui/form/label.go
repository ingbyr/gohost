package form

import (
	tea "github.com/charmbracelet/bubbletea"
)

var _ ItemModel = (*Label)(nil)

func NewLabel(label string) *Label {
	commonItem := NewCommonItem()
	commonItem.SetFocusable(false)
	return &Label{
		CommonItem: commonItem,
		label:      label,
	}
}

type Label struct {
	*CommonItem
	label         string
	width, height int
}

func (l *Label) Init() tea.Cmd {
	return nil
}

func (l *Label) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.width, l.height = msg.Width, msg.Height
	}
	return l, nil
}

func (l *Label) View() string {
	return l.label
}
