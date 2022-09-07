package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/config"
	"gohost/tui/styles"
	"strconv"
	"strings"
)

type sessionState int

const (
	treeViewState = iota
	editorViewState
	nodeViewState
	lastState
)

var cfg = config.Instance()

type Model struct {
	state          sessionState
	helpView       *HelpView
	groupView      *TreeView
	editorView     *EditorView
	nodeView       *NodeView
	reservedHeight int
	quitting       bool
}

func NewModel() (*Model, error) {
	model := &Model{
		state:          nodeViewState,
		reservedHeight: 6,
	}
	model.helpView = NewHelpView(model)
	model.groupView = NewTreeView(model)
	model.editorView = NewTextView(model)
	model.nodeView = NewNodeView(model)
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
	m.updateView(msg, &cmds, m.groupView)
	m.updateView(msg, &cmds, m.editorView)
	m.updateView(msg, &cmds, m.nodeView)
	m.updateView(msg, &cmds, m.helpView)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var b strings.Builder
	switch m.state {
	case treeViewState:
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			styles.FocusedView.Render(m.groupView.View()),
			styles.DefaultView.Render(m.editorView.View())))
	case editorViewState:
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			styles.DefaultView.Render(m.groupView.View()),
			styles.FocusedView.Render(m.editorView.View())))
	case nodeViewState:
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
			styles.DefaultView.Render(m.groupView.View()),
			styles.FocusedView.Render(m.nodeView.View())))
	}

	v := lipgloss.JoinVertical(lipgloss.Left, b.String(), m.helpView.View())
	return v
}

func (m *Model) updateView(msg tea.Msg, cmds *[]tea.Cmd, view tea.Model) {
	_, cmd := view.Update(msg)
	*cmds = append(*cmds, cmd)
}

func (m *Model) log(msg string) {
	m.helpView.debug = msg
}

func (m *Model) switchNextState() sessionState {
	m.switchState((m.state + 1) % lastState)
	m.log("state:" + strconv.Itoa(int(m.state)))
	return m.state
}

func (m *Model) switchState(state sessionState) {
	if state == editorViewState {
		m.editorView.Focus()
	} else {
		m.editorView.Blur()
	}
	m.state = state
}
