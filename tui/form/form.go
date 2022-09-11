package form

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/log"
	"gohost/tui/keys"
	"gohost/tui/styles"
	"strings"
)

func New() *Form {
	return &Form{
		Items:              make([]Item, 0),
		ItemFocusedStyle:   styles.None,
		ItemUnfocusedStyle: styles.None,
		Spacing:            0,
		preFocus:           0,
		focus:              0,
	}
}

type Form struct {
	Items              []Item
	ItemFocusedStyle   lipgloss.Style
	ItemUnfocusedStyle lipgloss.Style
	Spacing            int
	preFocus           int
	focus              int
	width              int
	height             int
}

func (v *Form) Init() tea.Cmd {
	return nil
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
	if len(v.Items) == 0 {
		return ""
	}
	h := v.Items[0].Height()
	if h > v.height {
		return ""
	}
	str := lipgloss.JoinVertical(lipgloss.Left, v.Items[0].View())

	for i := 1; i < len(v.Items); i++ {
		w := v.Items[i]
		if i == len(v.Items) - 1 {
			h += w.Height()
		} else {
			h += w.Height() + v.Spacing
		}
		if h >= v.height {
			return str
		}
		str += strings.Repeat(cfg.LineBreak, v.Spacing)
		str = lipgloss.JoinVertical(lipgloss.Left, str, v.Items[i].View())
		//log.Debug(fmt.Sprintf("cur h %d, view h %d", lipgloss.Height(str), v.height))
	}
	return str
}

func (v *Form) SetSize(width, height int) {
	v.width = width
	v.height = height
	remain := v.height - len(v.Items) * v.Spacing + 1
	height = v.height / len(v.Items)
	for i := 0; i < len(v.Items); i++ {
		w := v.Items[i]
		w.SetWidth(width)
		if i == len(v.Items)-1 {
			w.SetHeight(remain)
			log.Debug(fmt.Sprintf("base view w %d h %d", width, w.Height()))
		} else {
			w.SetHeight(height)
			remain -= w.Height()
			log.Debug(fmt.Sprintf("base view w %d h %d", width, w.Height()))
		}
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
