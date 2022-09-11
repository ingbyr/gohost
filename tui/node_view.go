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

	// Node type choices
	nodeTypeChoices := form.NewChoice([]list.DefaultItem{GroupNode, LocalHost, RemoteHost})
	nodeTypeChoices.Spacing = 1
	nodeTypeChoices.ShowMorePlaceHold = false

	// Confirm button
	confirmButton := form.NewButton("[ Confirm ]")
	confirmButton.OnClick = func() {
		log.Debug(fmt.Sprintf("name %s, desc %s, url %s, choice %s",
			nameTextInput.Value(), descTextInput.Value(), urlTextInput.Value(), nodeTypeChoices.SelectedItem()))
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
	nodeForm.AddItem(urlTextInput)
	nodeForm.AddItem(nodeTypeChoices)
	nodeForm.AddItem(confirmButton)

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
		log.Debug(fmt.Sprintf("node view w %d h %d", m.Width, m.Height))
	case tea.KeyMsg:
		if v.model.state == nodeViewState {
			return v.Form.Update(msg)
		}
	}
	_, cmd = v.Form.Update(msg)
	cmds = append(cmds, cmd)
	return v, tea.Batch(cmds...)
}

func (v *NodeView) View() string {
	return v.Form.View()
}
