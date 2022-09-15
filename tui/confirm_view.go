package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"gohost/tui/form"
	"gohost/tui/styles"
)

type ConfirmView struct {
	model *Model
	*form.Form
	tipLabel      *form.Label
	confirmButton *form.Button
	cancelButton  *form.Button
	width, height int
}

func NewConfirmView(model *Model) *ConfirmView {
	confirmForm := form.New()
	confirmForm.SetItemFocusedStyle(styles.FocusedFormItem)
	confirmForm.SetItemUnfocusedStyle(styles.UnfocusedFormItem)

	tipLabel := form.NewLabel("Default label")
	confirmForm.AddItem(tipLabel)

	confirmButton := form.NewButton("Confirm")
	confirmButton.OnClick = func() tea.Cmd {
		confirmForm.AddItem(form.NewLabel("click confirm button"))
		return nil
	}
	confirmForm.AddItem(confirmButton)

	cancelButton := form.NewButton("Cancel")
	cancelButton.OnClick = func() tea.Cmd {
		confirmForm.AddItem(form.NewLabel("click cancle button"))
		return nil
	}
	confirmForm.AddItem(cancelButton)

	confirmForm.FocusAvailableFirstItem()
	return &ConfirmView{
		model:         model,
		Form:          confirmForm,
		tipLabel:      tipLabel,
		confirmButton: confirmButton,
		cancelButton:  cancelButton,
	}
}

func (v *ConfirmView) Init() tea.Cmd {
	return nil
}

func (v *ConfirmView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.width, v.height = msg.Width, msg.Height
	case tea.KeyMsg:

	}
	cmds = append(cmds, cmd)
	_, cmd = v.Form.Update(msg)
	cmds = append(cmds, cmd)
	return v, tea.Batch(cmds...)
}

func (v *ConfirmView) Reset(tip string, confirmOnClick, cancelOnClick func() tea.Cmd) {
	v.tipLabel.Text = tip
	v.confirmButton.OnClick = confirmOnClick
	v.cancelButton.OnClick = cancelOnClick
}
