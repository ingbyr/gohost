package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type HelpView struct {
	view        help.Model
	debug       string
	enableDebug bool
}

func NewHelpView(model *Model) *HelpView {
	return &HelpView{
		view:        help.New(),
		enableDebug: true,
	}
}

func (h *HelpView) Init() tea.Cmd {
	return nil
}

func (h *HelpView) Update(msg tea.Msg) []tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.view.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Help):
			h.view.ShowAll = !h.view.ShowAll
		}
	}
	return nil
}

func (h *HelpView) View() string {
	var helper string
	helper += h.view.View(keys)
	if h.enableDebug {
		helper += "\nDebug: " + h.debug
	}
	return helper
}

func (h *HelpView) Width() int {
	return h.view.Width
}
