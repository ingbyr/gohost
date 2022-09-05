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
	docStyle          = lipgloss.NewStyle().Margin(1, 2)
	modelStyle        = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().Padding(1, 2).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
)

type Model struct {
	keys      keyMaps
	state     sessionState
	help      help.Model
	groupList list.Model
	hostList  list.Model
	quitting  bool
}

func NewModel() (*Model, error) {
	groupService := group.Service()
	groups, err := groupService.LoadGroups()
	if err != nil {
		return nil, err
	}
	groupItems := make([]list.Item, len(groups))
	for i := range groups {
		groupItems[i] = groups[i]
	}
	groupService.BuildTree(groups)

	keys := newKeys()
	m := &Model{
		keys:      keys,
		groupList: list.New(groupItems, list.NewDefaultDelegate(), 0, 0),
		help:      help.New(),
	}

	return m, nil
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
			cmd = tea.Quit
			cmds = append(cmds, cmd)
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
		v += m.groupList.View()
	}
	helpView := m.help.View(m.keys)
	v += helpView

	docStyle.Render(v)
	return v
}
