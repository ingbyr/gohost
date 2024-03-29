package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type HelpView struct {
	model       *Model
	helpModel   help.Model
	shortHelp   map[State][]key.Binding
	fullHelp    map[State][][]key.Binding
	enableDebug bool
}

func NewHelpView(model *Model) *HelpView {
	return &HelpView{
		model:       model,
		helpModel:   help.New(),
		shortHelp:   make(map[State][]key.Binding, 8),
		fullHelp:    make(map[State][][]key.Binding, 8),
		enableDebug: true,
	}
}

func (h *HelpView) Init() tea.Cmd {
	return nil
}

func (h *HelpView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.helpModel.Width = msg.Width
	}
	return h, nil
}

func (h *HelpView) View() string {
	return h.helpModel.View(h)
}

func (h *HelpView) ShortHelp() []key.Binding {
	return h.shortHelp[h.model.state]
}

func (h *HelpView) FullHelp() [][]key.Binding {
	return h.fullHelp[h.model.preState]
}

func (h *HelpView) Width() int {
	return h.helpModel.Width
}

func (h *HelpView) SetShortHelp(state State, kb []key.Binding) {
	h.shortHelp[state] = kb
}

func (h *HelpView) SetFullHelp(state State, kb [][]key.Binding) {
	h.fullHelp[state] = kb
}
