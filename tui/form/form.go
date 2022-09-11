package form

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/tui/keys"
	"gohost/tui/styles"
	"strings"
)

func New() *Form {
	return &Form{
		Items:              make([]Item, 0),
		ItemFocusedStyle:   styles.None,
		ItemUnfocusedStyle: styles.None,
		MorePlaceHold:      "...",
		Spacing:            0,
		viewport:           viewport.New(0, 0),
		preFocus:           0,
		focus:              0,
	}
}

type Form struct {
	Items              []Item
	ItemFocusedStyle   lipgloss.Style
	ItemUnfocusedStyle lipgloss.Style
	MorePlaceHold      string
	Spacing            int
	viewport           viewport.Model
	preFocus           int
	focus              int
	width              int
	height             int
}

func (v *Form) Init() tea.Cmd {
	return v.viewport.Init()
}

func (v *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.SetSize(m.Width, m.Height)
		return v, nil
	case tea.KeyMsg:
		focusedItem := v.Items[v.focus]
		switch {
		case key.Matches(m, keys.Up):
			intercepted := focusedItem.InterceptKey(m)
			_, cmd = focusedItem.Update(msg)
			if intercepted {
				return v, cmd
			}
			cmds = append(cmds, v.focusPreItem()...)
			return v, cmd
		case key.Matches(m, keys.Down):
			intercepted := focusedItem.InterceptKey(m)
			_, cmd = focusedItem.Update(msg)
			if intercepted {
				return v, cmd
			}
			cmds = append(cmds, v.focusNextItem()...)
			return v, cmd
		case key.Matches(m, keys.Enter):
			intercepted := focusedItem.InterceptKey(m)
			_, cmd = focusedItem.Update(msg)
			if intercepted {
				return v, cmd
			}
			cmds = append(cmds, v.focusNextItem()...)
			return v, cmd
		}
	}
	for i := 0; i < len(v.Items); i++ {
		_, cmd = v.Items[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return v, tea.Batch(cmds...)
}

func (v *Form) View() string {
	var b strings.Builder
	for i := range v.Items {
		b.WriteString(v.Items[i].View())
		if i < len(v.Items)-1 {
			b.WriteString(strings.Repeat(cfg.LineBreak, v.Spacing+1))
		}
	}
	v.viewport.SetContent(b.String())
	return v.viewport.View()
}

func (v *Form) SetSize(width, height int) {
	v.width = width
	v.height = height
	for _, item := range v.Items {
		item.SetWidth(width)
		v.viewport.Width = width
		v.viewport.Height = height
	}
}

func (v *Form) AddItem(widget Item) {
	if widget == nil {
		return
	}
	widget.SetFocusedStyle(v.ItemFocusedStyle)
	widget.SetUnfocusedStyle(v.ItemUnfocusedStyle)
	v.Items = append(v.Items, widget)
}

func (v *Form) SetItemFocusedStyle(style lipgloss.Style) {
	v.ItemFocusedStyle = style
	for i := range v.Items {
		v.Items[i].SetFocusedStyle(style)
	}
}

func (v *Form) SetItemUnfocusedStyle(style lipgloss.Style) {
	v.ItemUnfocusedStyle = style
	for i := range v.Items {
		v.Items[i].SetUnfocusedStyle(style)
	}
}

func (v *Form) focusNextItem() []tea.Cmd {
	nextFocus := v.idxAfterFocusItem()
	if nextFocus == v.focus {
		return nil
	}
	return v.setFocusItem(nextFocus, FocusFirstMode)
}

func (v *Form) focusPreItem() []tea.Cmd {
	nextFocus := v.idxBeforeFocusItem()
	if nextFocus == v.focus {
		return nil
	}
	return v.setFocusItem(v.idxBeforeFocusItem(), FocusLastMode)
}

func (v *Form) idxAfterFocusItem() int {
	idx := v.focus + 1
	if idx >= len(v.Items) {
		idx = 0
	}
	return idx
}

func (v *Form) idxBeforeFocusItem() int {
	idx := v.focus - 1
	if idx < 0 {
		idx = len(v.Items) - 1
	}
	return idx
}

func (v *Form) setFocusItem(idx int, mode FocusMode) []tea.Cmd {
	v.preFocus = v.focus
	v.focus = idx
	return []tea.Cmd{
		v.Items[v.preFocus].Unfocus(),
		v.Items[v.focus].Focus(mode),
	}
}
