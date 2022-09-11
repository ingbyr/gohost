package form

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/log"
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
		SelectedPrefix:    "[v]",
		UnselectedPrefix:  "[ ]",
		MorePlaceHold:     "...",
		ShowMorePlaceHold: true,
		width:             0,
		height:            0,
		Spacing:           1,
		focused:           false,
		cursorIndex:       -1,
		selectedIndex:     -1,
	}
}

type Choices struct {
	items             []list.DefaultItem
	SelectedPrefix    string
	UnselectedPrefix  string
	MorePlaceHold     string
	ShowMorePlaceHold bool
	focused           bool
	focusedStyle      lipgloss.Style
	unfocusedStyle    lipgloss.Style
	width, height     int
	Spacing           int
	cursorIndex       int
	selectedIndex     int
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

func (c *Choices) Width() int {
	return c.width
}

func (c *Choices) Height() int {
	return c.height
}

func (c *Choices) Init() tea.Cmd {
	return nil
}

func (c *Choices) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		c.SetHeight(m.Height)
		c.SetWidth(m.Width)
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
	log.Debug(fmt.Sprintf("choice cursor %d selected %d", c.cursorIndex, c.selectedIndex))
	cmds = append(cmds, cmd)
	return c, tea.Batch(cmds...)
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

func (c *Choices) SetWidth(width int) {
	c.width = width
}

func (c *Choices) SetHeight(height int) {
	c.height = height
	ah := len(c.items) + len(c.items)*c.Spacing - 1
	if ah < height {
		c.height = ah
	} else {
		c.height = height
	}
}

func (c *Choices) itemTitle(idx int) string {
	if idx == c.selectedIndex {
		return c.SelectedPrefix + c.items[idx].Title()
	}
	return c.UnselectedPrefix + c.items[idx].Title()
}
