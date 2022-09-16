package form

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/tui/keys"
	"gohost/tui/styles"
	"strings"
)

type FocusMode int

const (
	FocusFirstMode FocusMode = iota
	FocusLastMode
)

func New() *Form {
	return &Form{
		Items:                 make([]ItemModel, 0),
		defaultFocusedStyle:   styles.None,
		defaultUnfocusedStyle: styles.None,
		MorePlaceHold:         "...",
		Spacing:               0,
		preFocus:              0,
		focus:                 0,
		width:                 0,
		height:                0,
	}
}

type Form struct {
	Items                 []ItemModel
	defaultFocusedStyle   lipgloss.Style
	defaultUnfocusedStyle lipgloss.Style
	MorePlaceHold         string
	Spacing               int
	preFocus              int
	focus                 int
	width, height         int
	// TODO add layout align
}

func (v *Form) Init() tea.Cmd {
	return nil
}

func (v *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.width = m.Width
		v.height = m.Height
		// TODO add layout for items to adjust size dynamically
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
		item := v.Items[i]
		if item.Hided() {
			continue
		}
		b.WriteString(v.Items[i].View())
		if i < len(v.Items)-1 {
			b.WriteString(strings.Repeat("\n", v.Spacing+1))
		}
	}
	return lipgloss.NewStyle().
		Width(v.width).
		Height(v.height).
		Align(lipgloss.Top).
		Render(b.String())
}

func (v *Form) AddItem(widget ItemModel) {
	if widget == nil {
		return
	}
	widget.SetFocusedStyle(v.defaultFocusedStyle)
	widget.SetUnfocusedStyle(v.defaultUnfocusedStyle)
	v.Items = append(v.Items, widget)
}

func (v *Form) SetDefaultFocusedStyle(style lipgloss.Style) {
	v.defaultFocusedStyle = style
	for i := range v.Items {
		v.Items[i].SetFocusedStyle(style)
	}
}

func (v *Form) SetDefaultUnfocusedStyle(style lipgloss.Style) {
	v.defaultUnfocusedStyle = style
	for i := range v.Items {
		v.Items[i].SetUnfocusedStyle(style)
	}
}

func (v *Form) FocusAvailableFirstItem() {
	for i, item := range v.Items {
		if item.Focusable() {
			item.Focus(FocusFirstMode)
			v.focus = i
			return
		}
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
	idx := v.focus
	for {
		idx++
		if idx >= len(v.Items) {
			idx = 0
		}
		if !v.Items[idx].Hided() && v.Items[idx].Focusable() {
			return idx
		}
		if idx == v.focus {
			return v.focus
		}
	}
}

func (v *Form) idxBeforeFocusItem() int {
	idx := v.focus
	for {
		idx--
		if idx < 0 {
			idx = len(v.Items) - 1
		}
		if !v.Items[idx].Hided() && v.Items[idx].Focusable() {
			return idx
		}
		if idx == v.focus {
			return v.focus
		}
	}
}

func (v *Form) setFocusItem(idx int, mode FocusMode) []tea.Cmd {
	v.preFocus = v.focus
	v.focus = idx
	return []tea.Cmd{
		v.Items[v.preFocus].Unfocus(),
		v.Items[v.focus].Focus(mode),
	}
}
