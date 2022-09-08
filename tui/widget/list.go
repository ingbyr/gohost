package widget

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var _ Widget = (*List)(nil)

func NewList(items []list.Item, delegate list.ItemDelegate, width, height int) *List {
	return &List{
		Model: list.New(items, delegate, width, height),
	}
}

type List struct {
	list.Model
}

func (l *List) Init() tea.Cmd {
	return nil
}

func (l *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	l.Model, cmd = l.Model.Update(msg)
	return l, cmd
}

func (l *List) Focus() tea.Cmd {
	return nil
}

func (l *List) Unfocus() tea.Cmd {
	return nil
}

func (l *List) HasNext() bool {
	return l.Model.Index() != len(l.Model.Items())-1
}

func (l *List) HasPre() bool {
	return l.Model.Index() != 0
}
