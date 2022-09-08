package widget

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var _ Widget = (*Choice)(nil)

func NewChoice(items []list.Item, delegate list.ItemDelegate, width, height int) *Choice {
	return &Choice{
		Model: list.New(items, delegate, width, height),
	}
}

type Choice struct {
	list.Model
	selectIndexAfterUnfocus int
}

func (l *Choice) Init() tea.Cmd {
	return nil
}

func (l *Choice) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		l.SetHeight(m.Height / 5)
		l.SetWidth(m.Width / 5)
	}
	l.Model, cmd = l.Model.Update(msg)
	cmds = append(cmds, cmd)
	return l, tea.Batch(cmds...)
}

func (l *Choice) Focus() tea.Cmd {
	l.Select(l.selectIndexAfterUnfocus)
	return nil
}

func (l *Choice) Unfocus() tea.Cmd {
	l.selectIndexAfterUnfocus = l.Index()
	return nil
}

func (l *Choice) HandleKeyUp() bool {
	return !(l.Model.Index() == 0)
}

func (l *Choice) HandleKeyDown() bool {
	return !(l.Model.Index() == len(l.Model.Items())-1)
}
