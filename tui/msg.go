package tui

import tea "github.com/charmbracelet/bubbletea"

type RefreshTreeViewItems struct {
}

type AppliedNewHostContent struct {
}

type ConfirmMessage struct {
	Message       string
	ConfirmAction func() tea.Cmd
}
