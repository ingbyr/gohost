package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/config"
	"gohost/tui/styles"
	"strconv"
)

type sessionState int

const (
	treeViewState = iota
	editorViewState
	nodeViewState
	helpViewState
	lastState
)

var cfg = config.Instance()

type Model struct {
	preState       sessionState
	state          sessionState
	logView        *LogView
	helpView       *HelpView
	treeView       *TreeView
	editorView     *EditorView
	nodeView       *NodeView
	reservedHeight int
	quitting       bool
}

func NewModel() (*Model, error) {
	model := &Model{
		state:          treeViewState,
		reservedHeight: 1,
	}
	model.logView = NewLogView(model)
	model.helpView = NewHelpView(model)
	model.treeView = NewTreeView(model)
	model.editorView = NewTextView(model)
	model.nodeView = NewNodeView(model)
	return model, nil
}

func (m *Model) Init() tea.Cmd {
	m.log(fmt.Sprintf("view h %d w %d", styles.DefaultView.GetHeight(), styles.DefaultView.GetWidth()))
	return tea.Batch(
		m.treeView.Init(),
		m.editorView.Init(),
		m.nodeView.Init(),
		m.logView.Init(),
		m.helpView.Init(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Overwrite size msg for each component
		m.resizeViews(msg, &cmds)
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Switch):
			m.switchNextState()
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, keys.Help):
			if m.state != helpViewState {
				m.switchState(helpViewState)
				m.helpView.helpView.ShowAll = true
			} else {
				m.switchState(m.preState)
				m.helpView.helpView.ShowAll = false
			}
		}
	}
	m.updateView(msg, &cmds, m.logView)
	m.updateView(msg, &cmds, m.treeView)
	m.updateView(msg, &cmds, m.editorView)
	m.updateView(msg, &cmds, m.nodeView)
	m.updateView(msg, &cmds, m.helpView)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var v string
	switch m.state {
	case treeViewState:
		v = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				styles.DefaultView.Render(m.logView.View()),
				styles.FocusedView.Render(m.treeView.View()),
				styles.DefaultView.Render(m.editorView.View()),
			),
			m.helpView.View())

	case editorViewState:
		v = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				styles.DefaultView.Render(m.logView.View()),
				styles.DefaultView.Render(m.treeView.View()),
				styles.FocusedView.Render(m.editorView.View()),
			),
			m.helpView.View())

	case nodeViewState:
		v = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.JoinHorizontal(lipgloss.Top,
				styles.DefaultView.Render(m.logView.View()),
				styles.DefaultView.Render(m.treeView.View()),
				styles.FocusedView.Render(m.nodeView.View()),
			),
			m.helpView.View(),
		)

	case helpViewState:
		v = m.helpView.View()
	}
	return v
}

func (m *Model) updateView(msg tea.Msg, cmds *[]tea.Cmd, view tea.Model) {
	_, cmd := view.Update(msg)
	*cmds = append(*cmds, cmd)
}

func (m *Model) log(msg string) {
	m.logView.InsertLog(msg)
}

func (m *Model) switchNextState() sessionState {
	m.switchState((m.state + 1) % lastState)
	m.log("state:" + strconv.Itoa(int(m.state)))
	return m.state
}

func (m *Model) switchState(state sessionState) {
	m.preState = m.state
	if state == editorViewState {
		m.editorView.Focus()
	} else {
		m.editorView.Blur()
	}
	m.state = state
}

func (m *Model) setShortHelp(state sessionState, kb []key.Binding) {
	m.helpView.SetShortHelp(state, kb)
}

func (m *Model) setFullHelp(state sessionState, kb [][]key.Binding) {
	m.helpView.SetFullHelp(state, kb)
}

func (m *Model) resizeViews(sizeMsg tea.WindowSizeMsg, cmds *[]tea.Cmd) {
	m.log(fmt.Sprintf("window w %d h %d", sizeMsg.Width, sizeMsg.Height))
	width := sizeMsg.Width / 3
	height := sizeMsg.Height - m.reservedHeight
	m.updateView(tea.WindowSizeMsg{Width: width, Height: height}, cmds, m.treeView)
	m.updateView(tea.WindowSizeMsg{Width: width, Height: height}, cmds, m.editorView)
	m.updateView(tea.WindowSizeMsg{Width: width, Height: height}, cmds, m.nodeView)
	m.updateView(tea.WindowSizeMsg{Width: width, Height: height}, cmds, m.logView)
	m.updateView(tea.WindowSizeMsg{Width: sizeMsg.Width, Height: 1}, cmds, m.helpView)
	//m.log(fmt.Sprintf("w: %d h: %d", sizeMsg.Width, sizeMsg.Height))
}
