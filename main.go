package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"gohost/tui"
)

func main() {
	//if err := cmd.Execute(); err != nil {
	//	panic(err)
	//}

	m, err := tui.NewModel()
	if err != nil {
		panic(err)
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		panic(err)
	}
}
