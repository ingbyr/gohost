package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/group"
	"io"
	"strings"
)

type groupItem struct {
	*group.Node
	isFold bool
}

func WrapGroupNode(groupNode any) list.Item {
	node, ok := groupNode.(*group.Node)
	if !ok {
		return nil
	}

	return &groupItem{
		Node:   node,
		isFold: true,
	}
}

type groupItemDelegate struct {
}

func (d groupItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	groupNode, ok := item.(*groupItem)
	if !ok {
		return
	}
	var str string
	spaces := strings.Repeat(" ", groupNode.Depth)
	if m.Index() == index {
		str = fmt.Sprintf("> %s%d. %s", spaces, index+1, item.FilterValue())
	} else {
		str = fmt.Sprintf("  %s%d. %s", spaces, index+1, groupNode.Name)
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
