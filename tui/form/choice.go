package form

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/tui/keys"
	"strings"
)

type FocusMode int

const (
	FocusFirstMode FocusMode = iota
	FocusLastMode
)

var _ Item = (*Choices)(nil)

func NewChoice(items []list.DefaultItem) *Choices {
	return &Choices{
		items:             items,
		SelectedPrefix:    "(v) ",
		UnselectedPrefix:  "( ) ",
		MorePlaceHold:     "...",
		ShowMorePlaceHold: true,
		Spacing:           1,
		focused:           false,
		cursorIndex:       -1,
		selectedIndex:     -1,
	}
}

type Choices struct {
	SelectedPrefix    string
	UnselectedPrefix  string
	MorePlaceHold     string
	ShowMorePlaceHold bool
	HideFunc          HideCondition
	items             []list.DefaultItem
	focused           bool
	focusedStyle      lipgloss.Style
	unfocusedStyle    lipgloss.Style
	Spacing           int
	cursorIndex       int
	selectedIndex     int
}

func (c *Choices) Hide() bool {
	if c.HideFunc == nil {
		return false
	}
	return c.HideFunc()
}

func (c *Choices) SetFocusedStyle(style lipgloss.Style) {
	c.focusedStyle = style
}

func (c *Choices) SetUnfocusedStyle(style lipgloss.Style) {
	c.unfocusedStyle = style
}

func (c *Choices) Items() []list.DefaultItem {
	return c.items
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
			b.WriteString(strings.Repeat(cfg.LineBreak, c.Spacing+1))
		}
	}
	return b.String()
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
			}
		}
	}
	//log.Debug(fmt.Sprintf("choice cursor %d selected %d", c.cursorIndex, c.selectedIndex))
	return c, nil
}

func (c *Choices) Focus(mode FocusMode) tea.Cmd {
	c.focused = true
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
	c.focused = false
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
