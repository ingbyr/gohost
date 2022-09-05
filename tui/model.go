package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/group"
)

type sessionState uint

const (
	groupView = iota
	editorView
	sysHostView
)

var (
	docStyle          = lipgloss.NewStyle().Margin(1, 1)
	modelStyle        = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
)

type Model struct {
	keys          keyMaps
	state         sessionState
	help          help.Model
	groupList     list.Model
	selectedGroup *group.Node
	hostList      list.Model
	quitting      bool

	groupService *group.Service
}

func NewModel() (*Model, error) {
	// Get group service
	groupService := group.GetService()
	if err := groupService.Load(); err != nil {
		return nil, err
	}
	groups := wrapListItems(groupService.Tree())

	// Create group list view
	groupList := list.New(groups, list.NewDefaultDelegate(), 0, 0)
	// TODO add remaining help key
	groupList.Title = "Groups"
	groupList.SetShowHelp(false)
	return &Model{
		keys:         newKeys(),
		groupList:    groupList,
		help:         help.New(),
		groupService: groupService,
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
		h, v := docStyle.GetFrameSize()
		m.groupList.SetSize(msg.Width-h, msg.Height-v)
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
		case key.Matches(msg, m.keys.Down):
		case key.Matches(msg, m.keys.Left):
		case key.Matches(msg, m.keys.Right):
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, m.keys.Enter):
			m.selectedGroup = m.groupList.SelectedItem().(*group.Node)
		}
	}

	m.groupList, cmd = m.groupList.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var v string
	switch m.state {
	case groupView:
		v += lipgloss.JoinHorizontal(lipgloss.Top, m.groupList.View())
	}
	helpView := lipgloss.JoinVertical(lipgloss.Bottom, docStyle.Render(m.help.View(m.keys)))
	v += helpView
	docStyle.Render(v)
	// Debug
	if m.selectedGroup != nil {
		v += lipgloss.JoinVertical(lipgloss.Top, m.selectedGroup.Name)
	}
	return v
}
