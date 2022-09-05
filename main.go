package main

import tea "github.com/charmbracelet/bubbletea"

func main() {
	//if err := cmd.Execute(); err != nil {
	//	panic(err)
	//}

	m := NewModel()
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		panic(err)
	}
}
