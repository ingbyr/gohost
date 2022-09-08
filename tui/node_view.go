package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/tui/styles"
	"gohost/tui/view"
	"gohost/tui/widget"
	"strings"
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
	// TODO get width and height from NewNodeView args
	nodeTypes := widget.NewChoice([]list.Item{GroupNode, LocalHost, RemoteHost}, list.NewDefaultDelegate(), 2, 20)

	nodeView := &NodeView{
		model:       model,
		BaseView:    view.New(),
		preFocusIdx: 0,
		focusIdx:    0,
		nodeTypes:   nodeTypes,
	}
	nodeView.WidgetStyle = styles.DefaultView
	//nodeView.AddWidget(nodeNameTextInput)
	//nodeView.AddWidget(descTextInput)
	nodeView.AddWidget(urlTextInput)
	nodeView.AddWidget(nodeTypes)

	return nodeView
}

func (v *NodeView) Init() tea.Cmd {
	return nil
}

func (v *NodeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.nodeTypes.SetHeight(m.Height - v.model.reservedHeight)
		v.nodeTypes.SetWidth(m.Width / 3 * 2)
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

	//v.nodeTypes, cmd = v.nodeTypes.Update(msg)
	//cmds = append(cmds, cmd)

	return v, tea.Batch(cmds...)
}

func (v *NodeView) View() string {
	var b strings.Builder
	b.WriteString(v.BaseView.View())
	//b.WriteString("\n")
	//b.WriteString(v.nodeTypes.View())
	return b.String()
}
