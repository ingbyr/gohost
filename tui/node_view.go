package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/log"
	"gohost/tui/form"
	"gohost/tui/keys"
	"gohost/tui/styles"
)

type NodeView struct {
	model *Model
	*form.Form
	//nodeTypeChoices *form.Choices
}

func NewNodeView(model *Model) *NodeView {
	// Text inputs
	nodeNameTextInput := form.NewTextInput()
	nodeNameTextInput.Prompt = "Name: "
	nodeNameTextInput.Focus(form.FocusFirstMode)

	descTextInput := form.NewTextInput()
	descTextInput.Prompt = "Description: "

	urlTextInput := form.NewTextInput()
	urlTextInput.Prompt = "Url: "

	// Node type choices
	nodeTypeChoices := form.NewChoice([]list.DefaultItem{GroupNode, LocalHost, RemoteHost})
	nodeTypeChoices.Spacing = 1
	nodeTypeChoices.ShowMorePlaceHold = false

	nodeForm := &NodeView{
		model: model,
		Form:  form.New(),
		//nodeTypeChoices: nodeTypeChoices,
	}
	nodeForm.Spacing = 1
	nodeForm.SetItemFocusedStyle(styles.FocusedFormItem)
	nodeForm.SetItemUnfocusedStyle(styles.UnfocusedFormItem)
	nodeForm.AddItem(nodeNameTextInput)
	nodeForm.AddItem(descTextInput)
	nodeForm.AddItem(urlTextInput)
	nodeForm.AddItem(nodeTypeChoices)

	return nodeForm
}

func (v *NodeView) Init() tea.Cmd {
	v.model.setShortHelp(nodeViewState, keys.Arrows())
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
			_, cmd = v.Form.Update(msg)
		} else {
			return nil, nil
		}
	}
	cmds = append(cmds, cmd)
	return v, tea.Batch(cmds...)
}

func (v *NodeView) View() string {
	return v.Form.View()
}
