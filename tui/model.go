package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState uint

const (
	groupViewState = iota
	editorViewState
	sysHostViewState
)

var (
	docStyle          = lipgloss.NewStyle().Margin(1, 1)
	modelStyle        = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
	keys              = newKeys()
)

type Model struct {
	state sessionState
	help  help.Model
	groupView *GroupView
	hostList      list.Model
	quitting      bool
}

func NewModel() (*Model, error) {
	groupView, err := NewGroupView()
	if err != nil {
		return nil, err
	}
	return &Model{
		help:      help.New(),
		groupView: groupView,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
		case key.Matches(msg, keys.Down):
		case key.Matches(msg, keys.Left):
		case key.Matches(msg, keys.Right):
		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			cmds = append(cmds, tea.Quit)
		}
	}
	m.groupView.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var v string
	switch m.state {
	case groupViewState:
		v += lipgloss.JoinHorizontal(lipgloss.Top, m.groupView.View())
	}
	helpView := lipgloss.JoinVertical(lipgloss.Bottom, docStyle.Render(m.help.View(keys)))
	v += helpView
	docStyle.Render(v)
	return v
}
