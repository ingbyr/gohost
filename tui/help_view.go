package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
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

func (h *HelpView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.view.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Help):
			h.view.ShowAll = !h.view.ShowAll
		}
	}
	return h, nil
}

func (h *HelpView) View() string {
	var b strings.Builder
	b.WriteString(h.view.View(keys))
	if h.enableDebug {
		b.WriteString(cfg.LineBreak)
		b.WriteString("Debug: ")
		b.WriteString(h.debug)
	}
	return b.String()
}

func (h *HelpView) Width() int {
	return h.view.Width
}
