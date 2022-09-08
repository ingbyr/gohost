package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"gohost/tui/widget"
	"strings"
)

var _ View = (*BaseView)(nil)

func New() *BaseView {
	return &BaseView{
		widgets:  make([]widget.Widget, 0),
		preFocus: 0,
		focus:    0,
	}
}

type BaseView struct {
	widgets  []widget.Widget
	preFocus int
	focus    int
}

func (v *BaseView) Init() tea.Cmd {
	panic("implement me")
}

func (v *BaseView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := v.widgets[v.focus].Update(msg)
	return v, cmd
}

func (v *BaseView) View() string {
	var b strings.Builder
	for i := 0; i < len(v.widgets); i++ {
		b.WriteString(v.widgets[i].View())
		b.WriteString("\n")
	}
	return b.String()
}

func (v *BaseView) AddWidget(widget widget.Widget) {
	if widget == nil {
		return
	}
	v.widgets = append(v.widgets, widget)
}

func (v *BaseView) FocusNextWidget() []tea.Cmd {
	return v.setFocusWidget(v.idxAfterFocusWidget())
}

func (v *BaseView) FocusPreWidget() []tea.Cmd {
	return v.setFocusWidget(v.idxBeforeFocusWidget())
}

func (v *BaseView) idxAfterFocusWidget() int {
	if v.widgets[v.focus].HasNext() {
		return v.focus
	}
	idx := v.focus + 1
	if idx >= len(v.widgets) {
		idx = 0
	}
	return idx
}

func (v *BaseView) idxBeforeFocusWidget() int {
	if v.widgets[v.focus].HasPre() {
		return v.focus
	}
	idx := v.focus - 1
	if idx < 0 {
		idx = len(v.widgets) - 1
	}
	return idx
}

func (v *BaseView) setFocusWidget(idx int) []tea.Cmd {
	v.preFocus = v.focus
	v.focus = idx
	return []tea.Cmd{
		v.widgets[v.preFocus].Unfocus(),
		v.widgets[v.focus].Focus(),
	}
}
