package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/group"
	"io"
	"strings"
)

type groupItemDelegate struct {
}

func (d groupItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	node, ok := item.(*group.Node)
	if !ok {
		return
	}
	var str string
	spaces := strings.Repeat(" ", node.Depth)
	if m.Index() == index {
		str = fmt.Sprintf("> %s%d. %s", spaces, index+1, item.FilterValue())
	} else {
		str = fmt.Sprintf("  %s%d. %s", spaces, index+1, node.Name)
	}
	_, _ = fmt.Fprint(w, str)
}

func (d groupItemDelegate) Height() int {
	return 1
}

func (d groupItemDelegate) Spacing() int {
	return 0
}

func (d groupItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
