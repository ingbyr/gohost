package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/gohost"
	"gohost/util"
	"io"
	"strings"
)

// groupItemDelegate is item delegate for groupItem
type groupItemDelegate struct {
}

func (d groupItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	node, ok := item.(*gohost.Node[gohost.TreeNode])
	if !ok {
		return
	}
	var str string
	if m.Index() == index {
		str = "> "
	} else {
		str = "  "
	}
	spaces := strings.Repeat(" ", node.Depth)
	switch node := node.Data.(type) {
	case gohost.Group:
		str += fmt.Sprintf("%s[G] %d. %s", spaces, index, node.Name)
	case *gohost.LocalHost:
		str += fmt.Sprintf("%s[L] %d. %s", spaces, index, node.Name)
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

// GroupView is tui view for nodes tree
type GroupView struct {
	model         *Model
	groupList     list.Model
	selectedNode  *gohost.Node[gohost.TreeNode]
	selectedIndex int

	service *gohost.Service
}

func NewGroupView(model *Model) *GroupView {
	// Get nodes service
	service := gohost.GetService()
	service.Load()
	groups := util.WrapSlice[list.Item](service.Tree())

	// Create nodes list view
	groupList := list.New(groups, groupItemDelegate{}, 0, 0)
	// TODO add remaining help key
	groupList.Title = "Groups"
	groupList.SetShowHelp(false)

	return &GroupView{
		groupList:    groupList,
		selectedNode: nil,
		service:      service,
	}
}

func (v *GroupView) Init() tea.Cmd {
	return nil
}

func (v *GroupView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w, h := docStyle.GetFrameSize()
		v.groupList.SetSize(msg.Width-w, msg.Height-h)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			selectedItem := v.groupList.SelectedItem()
			if selectedItem != nil {
				v.selectedNode = selectedItem.(*gohost.Node[gohost.TreeNode])
				v.selectedIndex = v.groupList.Index()
				switch node := v.selectedNode.Data.(type) {
				case gohost.Group:
					cmds = v.onGroupNodeEnterClick(cmds)
					v.selectedNode.IsFolded = !v.selectedNode.IsFolded
				case gohost.Host:
					v.groupList.Title = "select host " + node.GetName()
				}
			}
		}
	}

	v.groupList, cmd = v.groupList.Update(msg)
	return v, tea.Batch(append(cmds, cmd)...)
}

func (v *GroupView) View() string {
	return v.groupList.View()
}

func (v *GroupView) onGroupNodeEnterClick(cmds []tea.Cmd) []tea.Cmd {
	if v.selectedNode.IsFolded {
		cmds = v.unfoldSelectedGroup(cmds)
	} else {
		v.foldSelectedGroup()
	}
	return cmds
}

func (v *GroupView) unfoldSelectedGroup(cmds []tea.Cmd) []tea.Cmd {
	subGroups := v.service.ChildNodes(v.selectedNode.GetID())
	idx := v.selectedIndex
	for i := range subGroups {
		idx++
		cmds = append(cmds, v.groupList.InsertItem(idx, subGroups[i]))
	}
	subHosts := v.service.LoadHostNodes(v.selectedNode.GetID())
	for i := range subHosts {
		idx++
		cmds = append(cmds, v.groupList.InsertItem(idx, subHosts[i]))
	}
	return cmds
}

func (v *GroupView) foldSelectedGroup() {
	items := v.groupList.Items()
	next := v.selectedIndex + 1
	for i := next; i < len(items); i++ {
		if items[next] == nil {
			break
		}
		node := items[next].(*gohost.Node[gohost.TreeNode])
		if node.Depth > v.selectedNode.Depth {
			v.groupList.RemoveItem(next)
		}
	}
}
