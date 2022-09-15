package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/tui/styles"
)

type ItemModel interface {
	tea.Model
	Item
}

type Item interface {
	Focus(mode FocusMode) tea.Cmd
	Unfocus() tea.Cmd
	Focusable() bool
	SetFocusable(focusable bool)
	SetFocusedStyle(style lipgloss.Style)
	SetUnfocusedStyle(style lipgloss.Style)
	Hided() bool
	SetHided(hided bool)
	InterceptKey(m tea.KeyMsg) bool
}

func NewCommonItem() *CommonItem {
	return &CommonItem{
		focusable:      true,
		focused:        false,
		focusedStyle:   styles.None,
		unfocusedStyle: styles.None,
		hided:          false,
	}
}

var _ Item = (*CommonItem)(nil)

type CommonItem struct {
	focusable      bool
	focused        bool
	focusedStyle   lipgloss.Style
	unfocusedStyle lipgloss.Style
	hided          bool
}

func (c *CommonItem) Focus(mode FocusMode) tea.Cmd {
	c.focused = true
	return nil
}

func (c *CommonItem) Unfocus() tea.Cmd {
	c.focused = false
	return nil
}

func (c *CommonItem) Focusable() bool {
	return c.focusable
}

func (c *CommonItem) SetFocusable(focusable bool) {
	c.focusable = focusable
}

func (c *CommonItem) InterceptKey(m tea.KeyMsg) bool {
	return false
}

func (c *CommonItem) SetFocusedStyle(style lipgloss.Style) {
	c.focusedStyle = style
}

func (c *CommonItem) SetUnfocusedStyle(style lipgloss.Style) {
	c.unfocusedStyle = style
}

func (c *CommonItem) Hided() bool {
	return c.hided
}

func (c *CommonItem) SetHided(hided bool) {
	c.hided = hided
}
