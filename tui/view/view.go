package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/tui/styles"
	"gohost/tui/widget"
)

type View interface {
	tea.Model
	AddWidget(widget widget.Widget)
	FocusNextWidget() []tea.Cmd
	FocusPreWidget() []tea.Cmd
}

var _ View = (*BaseView)(nil)

func New() *BaseView {
	return &BaseView{
		Widgets:     make([]widget.Widget, 0),
		WidgetStyle: styles.None,
		preFocus:    0,
		focus:       0,
	}
}

type BaseView struct {
	Widgets     []widget.Widget
	WidgetStyle lipgloss.Style
	preFocus    int
	focus       int
}

func (v *BaseView) Init() tea.Cmd {
	panic("implement me")
}

func (v *BaseView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		for i := 0; i < len(v.Widgets); i++ {

		}
	case tea.KeyMsg:
		switch m.String() {
		case "up", "down":
			_, cmd = v.Widgets[v.focus].Update(msg)
			return v, cmd
		}
	}
	for i := 0; i < len(v.Widgets); i++ {
		_, cmd = v.Widgets[v.focus].Update(msg)
		cmds = append(cmds, cmd)
	}
	return v, tea.Batch(cmds...)
}

func (v *BaseView) View() string {
	var str string
	for i := 0; i < len(v.Widgets); i++ {
		str = lipgloss.JoinVertical(lipgloss.Left, str, v.Widgets[i].View())
	}
	return str
}

func (v *BaseView) AddWidget(widget widget.Widget) {
	if widget == nil {
		return
	}
	v.Widgets = append(v.Widgets, widget)
}

func (v *BaseView) FocusNextWidget() []tea.Cmd {
	return v.setFocusWidget(v.idxAfterFocusWidget())
}

func (v *BaseView) FocusPreWidget() []tea.Cmd {
	return v.setFocusWidget(v.idxBeforeFocusWidget())
}

func (v *BaseView) idxAfterFocusWidget() int {
	if v.Widgets[v.focus].HandleKeyDown() {
		return v.focus
	}
	idx := v.focus + 1
	if idx >= len(v.Widgets) {
		idx = 0
	}
	return idx
}

func (v *BaseView) idxBeforeFocusWidget() int {
	if v.Widgets[v.focus].HandleKeyUp() {
		return v.focus
	}
	idx := v.focus - 1
	if idx < 0 {
		idx = len(v.Widgets) - 1
	}
	return idx
}

func (v *BaseView) setFocusWidget(idx int) []tea.Cmd {
	v.preFocus = v.focus
	v.focus = idx
	return []tea.Cmd{
		v.Widgets[v.preFocus].Unfocus(),
		v.Widgets[v.focus].Focus(),
	}
}
