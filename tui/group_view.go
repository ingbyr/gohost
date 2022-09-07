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
	node, ok := item.(*gohost.TreeNode[gohost.Node])
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
	switch node := node.Node.(type) {
	case gohost.Group:
		str += fmt.Sprintf("%s[G] %d. %s", spaces, index, node.Name)
	case gohost.Host:
		str += fmt.Sprintf("%s[L] %d. %s", spaces, index, node.GetName())
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
	selectedNode  *gohost.TreeNode[gohost.Node]
	selectedIndex int
	selectedGroup gohost.Group
	selectedHost  gohost.Host

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
		model:        model,
		groupList:    groupList,
		selectedNode: service.SysHostNode,
		service:      service,
	}
}

func (v *GroupView) Init() tea.Cmd {
	return nil
}

func (v *GroupView) Update(msg tea.Msg) []tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.groupList.SetHeight(m.Height - v.model.reservedHeight)
		v.groupList.SetWidth(m.Width / 3)
		v.model.helpView.debug = fmt.Sprintf("w %d h %d, w %d h %d", m.Width, m.Height, v.groupList.Width(), v.groupList.Height())
	case tea.KeyMsg:
		if v.model.state == groupViewState {
			switch {
			case key.Matches(m, keys.Enter):
				selectedItem := v.groupList.SelectedItem()
				if selectedItem != nil {
					v.selectedNode = selectedItem.(*gohost.TreeNode[gohost.Node])
					v.selectedIndex = v.groupList.Index()
					switch v.selectedNode.Node.(type) {
					case gohost.Group:
						v.onGroupNodeEnterClick(&cmds)
					case gohost.Host:
						v.onHostNodeSelected(&cmds)
					}
				}
			}
		} else {
			// Disable key
			msg = nil
		}
	}
	v.groupList, cmd = v.groupList.Update(msg)
	return append(cmds, cmd)
}

func (v *GroupView) View() string {
	return v.groupList.View()
}

func (v *GroupView) onGroupNodeEnterClick(cmds *[]tea.Cmd) {
	v.selectedGroup = v.selectedNode.Node.(gohost.Group)
	if v.selectedNode.IsFolded {
		v.unfoldSelectedGroup(cmds)
	} else {
		v.foldSelectedGroup()
	}
	v.selectedNode.IsFolded = !v.selectedNode.IsFolded
}

func (v *GroupView) unfoldSelectedGroup(cmds *[]tea.Cmd) {
	subGroups := v.service.ChildNodes(v.selectedNode.GetID())
	idx := v.selectedIndex
	for i := range subGroups {
		idx++
		*cmds = append(*cmds, v.groupList.InsertItem(idx, subGroups[i]))
	}
	subHosts := v.service.LoadHostNodes(v.selectedNode.GetID())
	for i := range subHosts {
		idx++
		*cmds = append(*cmds, v.groupList.InsertItem(idx, subHosts[i]))
	}
}

func (v *GroupView) foldSelectedGroup() {
	items := v.groupList.Items()
	next := v.selectedIndex + 1
	for i := next; i < len(items); i++ {
		if items[next] == nil {
			break
		}
		node := items[next].(*gohost.TreeNode[gohost.Node])
		if node.Depth > v.selectedNode.Depth {
			node.IsFolded = true
			v.groupList.RemoveItem(next)
		} else {
			break
		}
	}
}

func (v *GroupView) onHostNodeSelected(cmds *[]tea.Cmd) {
	v.selectedHost = v.selectedNode.Node.(gohost.Host)
	v.model.Log("select host: " + v.selectedHost.GetName())
	v.model.SwitchState(editorViewState)
	v.model.editorView.SetHost(v.selectedHost)
}
