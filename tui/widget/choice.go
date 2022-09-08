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

func (c *Choice) Init() tea.Cmd {
	return nil
}

func (c *Choice) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		c.SetHeight(m.Height / 5)
		c.SetWidth(m.Width / 5)
	}
	c.Model, cmd = c.Model.Update(msg)
	cmds = append(cmds, cmd)
	return c, tea.Batch(cmds...)
}

func (c *Choice) Focus() tea.Cmd {
	c.Select(c.selectIndexAfterUnfocus)
	return nil
}

func (c *Choice) Unfocus() tea.Cmd {
	c.selectIndexAfterUnfocus = c.Index()
	return nil
}

func (c *Choice) HandleKeyUp() bool {
	return !(c.Model.Index() == 0)
}

func (c *Choice) HandleKeyDown() bool {
	return !(c.Model.Index() == len(c.Model.Items())-1)
}
