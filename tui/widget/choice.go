package widget

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/log"
)

var _ Widget = (*Choice)(nil)

func NewChoice(items []list.Item) *Choice {
	delegate := list.NewDefaultDelegate()
	model := list.New(items, delegate, 0, 0)
	model.SetShowHelp(false)
	model.SetShowPagination(true)
	model.SetShowStatusBar(false)
	model.SetFilteringEnabled(false)
	model.SetShowTitle(false)
	return &Choice{
		Model: model,
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
		c.SetHeight(m.Height)
		c.SetWidth(m.Width)
		log.Debug(fmt.Sprintf("choice w %d h %d", c.Model.Width(), c.Model.Height()))
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

func (c *Choice) SetWidth(width int) {
	c.Model.SetWidth(width)
}

func (c *Choice) SetHeight(height int) int {
	c.Model.SetHeight(height)
	return c.Model.Height()
}
