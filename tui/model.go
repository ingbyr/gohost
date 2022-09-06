package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strconv"
)

type sessionState int

const (
	groupViewState = iota
	editorViewState
	sysHostViewState
	lastState
)

var (
	modelStyle        = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
	keys              = newKeys()
)

type Model struct {
	state          sessionState
	helpView       *HelpView
	groupView      *GroupView
	editorView     *EditorView
	reservedHeight int
	quitting       bool
}

func NewModel() (*Model, error) {
	model := &Model{
		state:          0,
		reservedHeight: 6,
	}
	model.helpView = NewHelpView(model)
	model.groupView = NewGroupView(model)
	model.editorView = NewTextView(model)
	return model, nil
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.groupView.Init(),
		m.editorView.Init(),
		m.helpView.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Switch):
			m.switchNextState()
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			cmds = append(cmds, tea.Quit)
		}
	}
	cmds = append(cmds, m.groupView.Update(msg)...)
	cmds = append(cmds, m.editorView.Update(msg)...)
	cmds = append(cmds, m.helpView.Update(msg)...)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var str string
	switch m.state {
	case groupViewState:
		str = lipgloss.JoinHorizontal(lipgloss.Top,
			focusedModelStyle.Render(m.groupView.View()),
			modelStyle.Render(m.editorView.View()))
	case editorViewState:
		str = lipgloss.JoinHorizontal(lipgloss.Top,
			modelStyle.Render(m.groupView.View()),
			focusedModelStyle.Render(m.editorView.View()))
	}
	str = lipgloss.JoinVertical(lipgloss.Left, str, m.helpView.View())
	//str += "\n" + m.helpView.View()
	return str
}

func (m *Model) Log(msg string) {
	m.helpView.debug = msg
}

func (m *Model) switchNextState() sessionState {
	m.SwitchState((m.state + 1) % lastState)
	m.Log("state:" + strconv.Itoa(int(m.state)))
	return m.state
}

func (m *Model) SwitchState(state sessionState) {
	if state == editorViewState {
		m.editorView.Focus()
	} else {
		m.editorView.Blur()
	}
	m.state = state
}
