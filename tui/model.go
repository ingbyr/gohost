package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
	"strings"
)

type sessionState uint

const (
	groupViewState = iota
	editorViewState
	sysHostViewState
)

var (
	modelStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).BorderStyle(lipgloss.HiddenBorder())
	//focusedModelStyle = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
	keys = newKeys()
)

type Model struct {
	state     sessionState
	helpView  *HelpView
	groupView *GroupView
	textView  *TextView
	quitting  bool
}

func NewModel() (*Model, error) {
	model := &Model{
		state: groupViewState,
	}
	model.helpView = NewHelpView(model)
	model.groupView = NewGroupView(model)
	model.textView = NewTextView(model)
	return model, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			cmds = append(cmds, tea.Quit)
		}
	}
	cmds = append(cmds, m.groupView.Update(msg)...)
	cmds = append(cmds, m.textView.Update(msg)...)
	cmds = append(cmds, m.helpView.Update(msg)...)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var str string
	switch m.state {
	case groupViewState:
		str += lipgloss.JoinHorizontal(lipgloss.Top,
			modelStyle.Render(m.groupView.View()),
			modelStyle.Render(m.textView.View()))
	}
	//str += "\n" + m.debug
	helperStr := m.helpView.View()
	helperHeight := strings.Count(helperStr, "\n")
	m.helpView.debug = strconv.Itoa(helperHeight)
	str = lipgloss.JoinVertical(lipgloss.Left, str, m.helpView.View())
	return str
}
