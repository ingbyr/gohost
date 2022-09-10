package form

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/log"
	"gohost/tui/keys"
	"gohost/tui/styles"
)

type Form interface {
	tea.Model
	AddWidget(widget Widget)
	FocusNextWidget() []tea.Cmd
	FocusPreWidget() []tea.Cmd
	SetSize(width, height int)
}

var _ Form = (*BaseForm)(nil)

func New() *BaseForm {
	return &BaseForm{
		Widgets:     make([]Widget, 0),
		WidgetStyle: styles.None,
		preFocus:    0,
		focus:       0,
	}
}

type BaseForm struct {
	Widgets     []Widget
	WidgetStyle lipgloss.Style
	preFocus    int
	focus       int
	width       int
	height      int
}

func (v *BaseForm) Init() tea.Cmd {
	panic("implement me")
}

func (v *BaseForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.SetSize(m.Width, m.Height)
		return v, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(m, keys.Up), key.Matches(m, keys.Down):
			_, cmd = v.Widgets[v.focus].Update(msg)
			return v, cmd
		}
	}
	for i := 0; i < len(v.Widgets); i++ {
		_, cmd = v.Widgets[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return v, tea.Batch(cmds...)
}

func (v *BaseForm) View() string {
	var str string
	for i := 0; i < len(v.Widgets); i++ {
		w := v.Widgets[i]
		if w.Height()+lipgloss.Height(str) > v.height {
			return str
		}
		if i == 0 {
			str = lipgloss.JoinVertical(lipgloss.Left, w.View())
			continue
		}
		str = lipgloss.JoinVertical(lipgloss.Left, str, v.Widgets[i].View())
		//log.Debug(fmt.Sprintf("cur h %d, view h %d", lipgloss.Height(str), v.height))
	}
	return str
}
func (v *BaseForm) SetSize(width, height int) {
	v.width = width
	v.height = height
	remain := v.height
	height = v.height / len(v.Widgets)
	for i := 0; i < len(v.Widgets); i++ {
		w := v.Widgets[i]
		w.SetWidth(width)
		if i == len(v.Widgets)-1 {
			w.SetHeight(remain)
			log.Debug(fmt.Sprintf("base view w %d h %d", width, w.Height()))
		} else {
			w.SetHeight(height)
			remain -= w.Height()
			log.Debug(fmt.Sprintf("base view w %d h %d", width, w.Height()))
		}
	}
}

func (v *BaseForm) AddWidget(widget Widget) {
	if widget == nil {
		return
	}
	v.Widgets = append(v.Widgets, widget)
}

func (v *BaseForm) FocusNextWidget() []tea.Cmd {
	nextFocus := v.idxAfterFocusWidget()
	if nextFocus == v.focus {
		return nil
	}
	return v.setFocusWidget(nextFocus, FocusFirstMode)
}

func (v *BaseForm) FocusPreWidget() []tea.Cmd {
	nextFocus := v.idxBeforeFocusWidget()
	if nextFocus == v.focus {
		return nil
	}
	return v.setFocusWidget(v.idxBeforeFocusWidget(), FocusLastMode)
}

func (v *BaseForm) idxAfterFocusWidget() int {
	if v.Widgets[v.focus].HandleKeyDown() {
		return v.focus
	}
	idx := v.focus + 1
	if idx >= len(v.Widgets) {
		idx = 0
	}
	return idx
}

func (v *BaseForm) idxBeforeFocusWidget() int {
	if v.Widgets[v.focus].HandleKeyUp() {
		return v.focus
	}
	idx := v.focus - 1
	if idx < 0 {
		idx = len(v.Widgets) - 1
	}
	return idx
}

func (v *BaseForm) setFocusWidget(idx int, mode FocusMode) []tea.Cmd {
	v.preFocus = v.focus
	v.focus = idx
	return []tea.Cmd{
		v.Widgets[v.preFocus].Unfocus(),
		v.Widgets[v.focus].Focus(mode),
	}
}
