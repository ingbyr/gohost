package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"io"
)

type itemDelegate struct {
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	str := fmt.Sprintf("%d. %s", index+1, item.FilterValue())
	fmt.Fprint(w, str)
}

func (d itemDelegate) Height() int {
	return 1
}

func (d itemDelegate) Spacing() int {
	return 0
}

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
