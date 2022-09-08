package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type HelpView struct {
	model       *Model
	helpView    help.Model
	shortHelp   map[sessionState][]key.Binding
	fullHelp    map[sessionState][][]key.Binding
	enableDebug bool
}

func NewHelpView(model *Model) *HelpView {
	return &HelpView{
		model:       model,
		helpView:    help.New(),
		shortHelp:   make(map[sessionState][]key.Binding, 8),
		fullHelp:    make(map[sessionState][][]key.Binding, 8),
		enableDebug: true,
	}
}

func (h *HelpView) Init() tea.Cmd {
	return nil
}

func (h *HelpView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.helpView.Width = msg.Width
	}
	return h, nil
}

func (h *HelpView) View() string {
	return h.helpView.View(h)
}

func (h *HelpView) ShortHelp() []key.Binding {
	return h.shortHelp[h.model.state]
}

func (h *HelpView) FullHelp() [][]key.Binding {
	return h.fullHelp[h.model.preState]
}

func (h *HelpView) Width() int {
	return h.helpView.Width
}

func (h *HelpView) SetShortHelp(state sessionState, kb []key.Binding) {
	h.shortHelp[state] = kb
}

func (h *HelpView) SetFullHelp(state sessionState, kb [][]key.Binding) {
	h.fullHelp[state] = kb
}
