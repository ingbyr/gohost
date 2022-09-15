package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/gohost"
	"gohost/log"
	"gohost/tui/form"
	"gohost/tui/keys"
	"gohost/tui/styles"
)

type NodeView struct {
	model *Model
	*form.Form
	nameTextInput *form.TextInput
	descTextInput *form.TextInput
	urlTextInput  *form.TextInput
	typeChoices   *form.Choices
}

func NewNodeView(model *Model) *NodeView {
	// Text inputs
	nameTextInput := form.NewTextInput()
	nameTextInput.Prompt = "Name: "
	nameTextInput.Focus(form.FocusFirstMode)

	descTextInput := form.NewTextInput()
	descTextInput.Prompt = "Description: "

	urlTextInput := form.NewTextInput()
	urlTextInput.Prompt = "Url: "
	urlTextInput.SetHided(true)

	// Node type choices
	nodeTypeChoices := form.NewChoice([]list.DefaultItem{NodeGroup, NodeLocalHost, NodeRemoteHost})
	nodeTypeChoices.Spacing = 1
	nodeTypeChoices.ShowMorePlaceHold = false
	nodeTypeChoices.OnSelectedChoice = func(item list.DefaultItem) {
		if item == NodeRemoteHost {
			urlTextInput.SetHided(false)
		} else {
			urlTextInput.SetHided(true)
		}
	}

	// Confirm button
	confirmButton := form.NewButton("Confirm")
	confirmButton.OnClick = func() tea.Cmd {
		// Check inputs
		if nodeTypeChoices.SelectedItem() == nil {
			model.confirmView.Reset("Please select node type",
				func() tea.Cmd {
					model.setState(StateNodeView)
					return nil
				},
				func() tea.Cmd {
					model.setState(StateNodeView)
					return nil
				})
			model.setState(StateConfirmView)
			return nil
		}

		// Get parent group
		selectedNode := model.treeView.SelectedNode()
		var parent *gohost.TreeNode
		switch selectedNode.Node.(type) {
		case *gohost.Group:
			parent = selectedNode
		case gohost.Host:
			parent = selectedNode.Parent()
		}

		// Save node
		switch nodeTypeChoices.SelectedItem() {
		case NodeGroup:
			group := &gohost.Group{
				ParentID: parent.GetID(),
				Name:     nameTextInput.Value(),
				Desc:     descTextInput.Value(),
			}
			groupNode := gohost.NewTreeNode(group)
			groupNode.SetParent(parent)
			svc.SaveGroupNode(groupNode)
		case NodeLocalHost:
			localHost := &gohost.LocalHost{
				GroupID: parent.GetID(),
				Name:    nameTextInput.Value(),
				Content: nil,
				Desc:    descTextInput.Value(),
			}
			localHostNode := gohost.NewTreeNode(localHost)
			localHostNode.SetParent(parent)
			if err := svc.SaveHostNode(localHostNode); err != nil {
				panic(err)
			}
		case NodeRemoteHost:
		}

		// Go back to tree view state
		return func() tea.Msg {
			model.switchState(StateTreeView)
			return RefreshTreeViewItems{}
		}
	}

	nodeForm := &NodeView{
		model:         model,
		Form:          form.New(),
		nameTextInput: nameTextInput,
		descTextInput: descTextInput,
		urlTextInput:  urlTextInput,
		typeChoices:   nodeTypeChoices,
	}
	nodeForm.Spacing = 1
	nodeForm.SetItemFocusedStyle(styles.FocusedFormItem)
	nodeForm.SetItemUnfocusedStyle(styles.UnfocusedFormItem)
	nodeForm.AddItem(nameTextInput)
	nodeForm.AddItem(descTextInput)
	nodeForm.AddItem(nodeTypeChoices)
	nodeForm.AddItem(urlTextInput)
	nodeForm.AddItem(confirmButton)

	return nodeForm
}

func (v *NodeView) Init() tea.Cmd {
	v.model.setShortHelp(StateNodeView, keys.Arrows())
	return nil
}

func (v *NodeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		log.Debug(fmt.Sprintf("node view w %d h %d", m.Width, m.Height))
	case tea.KeyMsg:
		return v.Form.Update(msg)
	}
	_, cmd := v.Form.Update(msg)
	return v, cmd
}

func (v *NodeView) View() string {
	return v.Form.View()
}
