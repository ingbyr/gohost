package widget

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/log"
	"gohost/tui/styles"
)

type View interface {
	tea.Model
	AddWidget(widget Widget)
	FocusNextWidget() []tea.Cmd
	FocusPreWidget() []tea.Cmd
	SetSize(width, height int)
}

var _ View = (*BaseView)(nil)

func New() *BaseView {
	return &BaseView{
		Widgets:     make([]Widget, 0),
		WidgetStyle: styles.None,
		preFocus:    0,
		focus:       0,
	}
}

type BaseView struct {
	Widgets     []Widget
	WidgetStyle lipgloss.Style
	preFocus    int
	focus       int
	width       int
	height      int
}

func (v *BaseView) Init() tea.Cmd {
	panic("implement me")
}

func (v *BaseView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.SetSize(m.Width, m.Height)
		return v, nil
	case tea.KeyMsg:
		switch m.String() {
		case "up", "down":
			_, cmd = v.Widgets[v.focus].Update(msg)
			return v, cmd
		}
	default:
		for i := 0; i < len(v.Widgets); i++ {
			_, cmd = v.Widgets[v.focus].Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return v, tea.Batch(cmds...)
}

func (v *BaseView) View() string {
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
		log.Debug(fmt.Sprintf("cur h %d, view h %d", lipgloss.Height(str), v.height))
	}
	return str
}
func (v *BaseView) SetSize(width, height int) {
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

func (v *BaseView) AddWidget(widget Widget) {
	if widget == nil {
		return
	}
	v.Widgets = append(v.Widgets, widget)
}

func (v *BaseView) FocusNextWidget() []tea.Cmd {
	return v.setFocusWidget(v.idxAfterFocusWidget(), FocusFirstMode)
}

func (v *BaseView) FocusPreWidget() []tea.Cmd {
	return v.setFocusWidget(v.idxBeforeFocusWidget(), FocusLastMode)
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

func (v *BaseView) setFocusWidget(idx int, mode FocusMode) []tea.Cmd {
	v.preFocus = v.focus
	v.focus = idx
	return []tea.Cmd{
		v.Widgets[v.preFocus].Unfocus(),
		v.Widgets[v.focus].Focus(mode),
	}
}
