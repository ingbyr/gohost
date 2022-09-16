package form

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/tui/keys"
	"strings"
)

var _ ItemModel = (*Choices)(nil)

type OnSelectedChoice = func(item list.DefaultItem)

func NewChoice(items []list.DefaultItem) *Choices {
	return &Choices{
		CommonItem:        NewCommonItem(),
		items:             items,
		SelectedPrefix:    "(v) ",
		UnselectedPrefix:  "( ) ",
		MorePlaceHold:     "...",
		ShowMorePlaceHold: true,
		Spacing:           1,
		OnSelectedChoice:  nil,
		cursorIndex:       -1,
		selectedIndex:     -1,
	}
}

type Choices struct {
	*CommonItem
	SelectedPrefix    string
	UnselectedPrefix  string
	MorePlaceHold     string
	ShowMorePlaceHold bool
	OnSelectedChoice  OnSelectedChoice
	Spacing           int
	items             []list.DefaultItem
	cursorIndex       int
	selectedIndex     int
}

func (c *Choices) Items() []list.DefaultItem {
	return c.items
}

func (c *Choices) Init() tea.Cmd {
	return nil
}

func (c *Choices) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(m, keys.Up):
			if c.cursorIndex > 0 {
				c.cursorIndex--
			}
		case key.Matches(m, keys.Down):
			if c.cursorIndex < len(c.items)-1 {
				c.cursorIndex++
			}
		case key.Matches(m, keys.Enter):
			if c.cursorIndex >= 0 {
				c.selectedIndex = c.cursorIndex
				if c.OnSelectedChoice != nil {
					c.OnSelectedChoice(c.items[c.selectedIndex])
				}
			}
		}
	}
	//log.Debug(fmt.Sprintf("choice cursor %d selected %d", c.cursorIndex, c.selectedIndex))
	return c, nil
}

func (c *Choices) View() string {
	var b strings.Builder
	for i := range c.items {
		if i == c.cursorIndex {
			b.WriteString(c.focusedStyle.Render(c.itemTitle(i)))
		} else {
			b.WriteString(c.unfocusedStyle.Render(c.itemTitle(i)))
		}
		if i < len(c.items)-1 {
			b.WriteString(strings.Repeat("\n", c.Spacing+1))
		}
	}
	return b.String()
}

func (c *Choices) Focus(mode FocusMode) tea.Cmd {
	c.CommonItem.Focus(mode)
	if len(c.items) == 0 {
		return nil
	}
	if mode == FocusFirstMode {
		c.cursorIndex = 0
	} else if mode == FocusLastMode {
		c.cursorIndex = len(c.items) - 1
	}
	return nil
}

func (c *Choices) Unfocus() tea.Cmd {
	c.CommonItem.Unfocus()
	c.cursorIndex = -1
	return nil
}

func (c *Choices) InterceptKey(keyMsg tea.KeyMsg) bool {
	switch {
	case key.Matches(keyMsg, keys.Up):
		return !(c.cursorIndex == 0)
	case key.Matches(keyMsg, keys.Down):
		return !(c.cursorIndex == len(c.Items())-1)
	case key.Matches(keyMsg, keys.Enter):
		return true
	}
	return false
}

func (c *Choices) SelectedItem() list.DefaultItem {
	if c.selectedIndex == -1 {
		return nil
	}
	return c.items[c.selectedIndex]
}

func (c *Choices) itemTitle(idx int) string {
	if idx == c.selectedIndex {
		return c.SelectedPrefix + c.items[idx].Title()
	}
	return c.UnselectedPrefix + c.items[idx].Title()
}
