package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
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

	// Node type choices
	nodeTypeChoices := form.NewChoice([]list.DefaultItem{NodeGroup, NodeLocalHost, NodeRemoteHost})
	nodeTypeChoices.Spacing = 1
	nodeTypeChoices.ShowMorePlaceHold = false

	urlTextInput := form.NewTextInput()
	urlTextInput.Prompt = "Url: "
	urlTextInput.HideFunc = func() bool {
		return nodeTypeChoices.SelectedItem() != NodeRemoteHost
	}

	// Confirm button
	confirmButton := form.NewButton("[ Confirm ]")
	confirmButton.OnClick = func() tea.Cmd {
		log.Debug(fmt.Sprintf("name %s, desc %s, url %s, choice %s",
			nameTextInput.Value(), descTextInput.Value(), urlTextInput.Value(), nodeTypeChoices.SelectedItem()))
		// Check inputs
		if nodeTypeChoices.SelectedItem() == nil {
			log.Debug(fmt.Sprintf("no node type was selected"))
			return nil
		}

		// Get parent node
		selectedNode := model.treeView.SelectedNode()
		selectedNode.SetFolded(false)
		var parent *gohost.TreeNode
		switch selectedNode.Node.(type) {
		case *gohost.Group:
			parent = selectedNode
		case gohost.Host:
			parent = selectedNode.Parent()
		}

		var cmd tea.Cmd
		switch nodeTypeChoices.SelectedItem() {
		case NodeGroup:
			node := &gohost.Group{
				ParentID: parent.GetID(),
				Name:     nameTextInput.Value(),
				Desc:     descTextInput.Value(),
			}
			groupNode := gohost.NewTreeNode(node)
			groupNode.SetParent(parent)
			groupNode.SetDepth(parent.Depth() + 1)
			if err := svc.SaveGroupNode(groupNode); err != nil {
				panic(err)
			}
			cmd = model.treeView.RefreshTreeNodes()
		case NodeLocalHost:
		case NodeRemoteHost:
		}
		model.switchState(treeViewState)
		return cmd
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
	v.model.setShortHelp(nodeViewState, keys.Arrows())
	return nil
}

func (v *NodeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if v.model.state != nodeViewState {
		return v, nil
	}
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		log.Debug(fmt.Sprintf("node view w %d h %d", m.Width, m.Height))
	case tea.KeyMsg:
		if key.Matches(m, keys.Esc) {
			v.model.switchState(treeViewState)
		}
		return v.Form.Update(msg)
	}
	_, cmd := v.Form.Update(msg)
	return v, cmd
}

func (v *NodeView) View() string {
	return v.Form.View()
}
