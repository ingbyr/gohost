package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/log"
	"gohost/tui/view"
	"gohost/tui/widget"
)

var _ view.View = (*NodeView)(nil)

type NodeView struct {
	model *Model
	*view.BaseView
	preFocusIdx int
	focusIdx    int
	nodeTypes   *widget.Choice
}

func NewNodeView(model *Model) *NodeView {
	// Text inputs
	nodeNameTextInput := widget.NewTextInput()
	nodeNameTextInput.Prompt = "ID: "
	nodeNameTextInput.Focus()

	descTextInput := widget.NewTextInput()
	descTextInput.Prompt = "Description: "

	urlTextInput := widget.NewTextInput()
	urlTextInput.Prompt = "Url: "

	// Node type choices
	nodeTypes := widget.NewChoice([]list.Item{GroupNode, LocalHost, RemoteHost})
	nodeTypes.Title = "Type"

	nodeView := &NodeView{
		model:       model,
		BaseView:    view.New(),
		preFocusIdx: 0,
		focusIdx:    0,
		nodeTypes:   nodeTypes,
	}
	nodeView.WidgetStyle = lipgloss.NewStyle().PaddingBottom(1)
	nodeView.AddWidget(nodeNameTextInput)
	nodeView.AddWidget(descTextInput)
	nodeView.AddWidget(urlTextInput)
	nodeView.AddWidget(nodeTypes)

	return nodeView
}

func (v *NodeView) Init() tea.Cmd {
	v.model.setShortHelp(nodeViewState, keys.ArrowsHelp())
	return nil
}

func (v *NodeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.SetSize(m.Width, m.Height)
		log.Debug(fmt.Sprintf("node view w %d h %d", m.Width, m.Height))
	case tea.KeyMsg:
		if v.model.state == nodeViewState {
			switch {
			case key.Matches(m, keys.Enter, keys.Up, keys.Down):
				if key.Matches(m, keys.Enter, keys.Down) {
					cmds = append(cmds, v.FocusNextWidget()...)
				} else {
					cmds = append(cmds, v.FocusPreWidget()...)
				}
			}
		} else {
			return nil, tea.Batch(cmds...)
		}
	}

	_, cmd = v.BaseView.Update(msg)
	cmds = append(cmds, cmd)

	return v, tea.Batch(cmds...)
}

func (v *NodeView) View() string {
	return v.BaseView.View()
}
