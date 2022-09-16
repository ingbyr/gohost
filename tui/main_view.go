package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/config"
	"gohost/gohost"
	"gohost/log"
	"gohost/tui/keys"
	"gohost/tui/styles"
	"strconv"
)

type State int

const (
	StateTreeView = iota
	StateEditorView
	StateNodeView
	StateLastView
	StateHelpView
	StateConfirmView

	StateInit = StateTreeView
)

var (
	cfg = config.Instance()
	svc = gohost.GetService()
)

type Model struct {
	preState                State
	state                   State
	helpView                *HelpView
	confirmView             *ConfirmView
	treeView                *TreeView
	editorView              *EditorView
	nodeView                *NodeView
	width, height           int
	styleWidth, styleHeight int
	shortHelperHeight       int
	leftViewWidth           int
	rightViewWidth          int
}

func NewModel() (*Model, error) {
	styleWidth, styleHeight := styles.DefaultView.GetFrameSize()
	model := &Model{
		preState:          StateInit,
		state:             StateInit,
		styleWidth:        styleWidth * 2,
		styleHeight:       styleHeight,
		shortHelperHeight: 1,
	}
	model.helpView = NewHelpView(model)
	model.confirmView = NewConfirmView(model)
	model.treeView = NewTreeView(model)
	model.editorView = NewEditorView(model)
	model.nodeView = NewNodeView(model)
	return model, nil
}

func (m *Model) Init() tea.Cmd {
	log.Debug(fmt.Sprintf("style w %d h %d", m.styleWidth, m.styleHeight))
	return tea.Batch(
		m.helpView.Init(),
		m.confirmView.Init(),
		m.treeView.Init(),
		m.editorView.Init(),
		m.nodeView.Init(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeViews(msg, &cmds)
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Switch):
			m.switchNextState()
			m.helpView.helpModel.ShowAll = false
		case key.Matches(msg, keys.ForceQuit):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, keys.Esc):
			if m.state != StateInit {
				m.switchState(StateInit)
			} else {
				cmds = append(cmds, tea.Quit)
			}
		case key.Matches(msg, keys.Help):
			if m.state != StateHelpView {
				m.switchState(StateHelpView)
				m.helpView.helpModel.ShowAll = true
			} else {
				m.switchState(m.preState)
				m.helpView.helpModel.ShowAll = false
			}
		default:
			switch m.state {
			case StateTreeView:
				m.updateView(msg, &cmds, m.treeView)
			case StateEditorView:
				m.updateView(msg, &cmds, m.editorView)
			case StateNodeView:
				m.updateView(msg, &cmds, m.nodeView)
			case StateConfirmView:
				m.updateView(msg, &cmds, m.confirmView)
			}
		}
	default:
		m.updateView(msg, &cmds, m.editorView)
		m.updateView(msg, &cmds, m.nodeView)
		m.updateView(msg, &cmds, m.treeView)
		m.updateView(msg, &cmds, m.confirmView)
		m.updateView(msg, &cmds, m.helpView)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var v string
	switch m.state {
	case StateTreeView:
		v = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				styles.FocusedView.Render(m.treeView.View()),
				styles.DefaultView.Render(m.editorView.View()),
			),
			m.helpView.View())

	case StateEditorView:
		v = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				styles.DefaultView.Render(m.treeView.View()),
				styles.FocusedView.Render(m.editorView.View()),
			),
			m.helpView.View())

	case StateNodeView:
		v = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				styles.DefaultView.Render(m.treeView.View()),
				styles.FocusedView.Render(m.nodeView.View()),
			),
			m.helpView.View(),
		)
	case StateConfirmView:
		v = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				styles.DefaultView.Render(m.treeView.View()),
				styles.FocusedView.Render(m.confirmView.View()),
			),
			m.helpView.View(),
		)
	case StateHelpView:
		v = m.helpView.View()
	}
	return v
}

func (m *Model) updateView(msg tea.Msg, cmds *[]tea.Cmd, view tea.Model) {
	_, cmd := view.Update(msg)
	*cmds = append(*cmds, cmd)
}

func (m *Model) switchNextState() State {
	m.switchState((m.state + 1) % StateLastView)
	log.Debug("state:" + strconv.Itoa(int(m.state)))
	return m.state
}

func (m *Model) switchState(state State) {
	m.preState = m.state
	m.setState(state)
}

func (m *Model) setState(state State) {
	m.state = state
	m.triggerStateViewIfNecessary(state)
}

func (m *Model) triggerStateViewIfNecessary(state State) {
	if state == StateEditorView {
		m.editorView.Focus()
	} else {
		m.editorView.Blur()
	}
}
func (m *Model) setShortHelp(state State, kb []key.Binding) {
	m.helpView.SetShortHelp(state, kb)
}

func (m *Model) setFullHelp(state State, kb [][]key.Binding) {
	m.helpView.SetFullHelp(state, kb)
}

func (m *Model) resizeViews(sizeMsg tea.WindowSizeMsg, cmds *[]tea.Cmd) {
	log.Debug(fmt.Sprintf("window w %d h %d", sizeMsg.Width, sizeMsg.Height))
	width := sizeMsg.Width - m.styleWidth
	m.leftViewWidth = width / 4
	m.rightViewWidth = width - m.leftViewWidth
	height := sizeMsg.Height - m.styleHeight - m.shortHelperHeight
	log.Debug(fmt.Sprintf("left w %d right w %d h %d", m.leftViewWidth, m.rightViewWidth, height))
	m.updateView(tea.WindowSizeMsg{Width: sizeMsg.Width, Height: 1}, cmds, m.helpView)
	m.updateView(tea.WindowSizeMsg{Width: m.leftViewWidth, Height: height}, cmds, m.treeView)
	m.updateView(tea.WindowSizeMsg{Width: m.rightViewWidth, Height: height}, cmds, m.confirmView)
	m.updateView(tea.WindowSizeMsg{Width: m.rightViewWidth, Height: height}, cmds, m.editorView)
	m.updateView(tea.WindowSizeMsg{Width: m.rightViewWidth, Height: height}, cmds, m.nodeView)
}
