package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/host"
	"io"
	"strings"
)

// groupItemDelegate is item delegate for groupItem
type groupItemDelegate struct {
}

func (d groupItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	node, ok := item.(*host.Node[host.TreeNode])
	if !ok {
		return
	}
	var str string
	spaces := strings.Repeat(" ", node.Depth)
	switch node := node.Data.(type) {
	case host.Group:
		if m.Index() == index {
			str = fmt.Sprintf("> %s%d. %s", spaces, index, node.Name)
		} else {
			str = fmt.Sprintf("  %s%d. %s", spaces, index, node.Name)
		}
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
	groupList     list.Model
	selectedGroup *host.Node[host.TreeNode]
	selectedIndex int

	groupService *host.Service
}

func NewGroupView() (*GroupView, error) {
	// Get nodes service
	groupService := host.GetService()
	if err := groupService.Load(); err != nil {
		return nil, err
	}
	groups := wrapListItems(groupService.Tree())

	// Create nodes list view
	groupList := list.New(groups, groupItemDelegate{}, 0, 0)
	// TODO add remaining help key
	groupList.Title = "Groups"
	groupList.SetShowHelp(false)

	return &GroupView{
		groupList:     groupList,
		selectedGroup: nil,
		groupService:  groupService,
	}, nil
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
				v.selectedGroup = selectedItem.(*host.Node[host.TreeNode])
				v.selectedIndex = v.groupList.Index()
				children := v.groupService.Children(v.selectedGroup.GetID())
				if v.selectedGroup.IsFold {
					for i := range children {
						cmds = append(cmds, v.groupList.InsertItem(v.selectedIndex+i+1, children[i]))
					}
				} else {
					v.foldSelectedGroup()
				}
				v.selectedGroup.IsFold = !v.selectedGroup.IsFold
			}
		}
	}

	v.groupList, cmd = v.groupList.Update(msg)
	return v, tea.Batch(append(cmds, cmd)...)
}

func (v *GroupView) View() string {
	return v.groupList.View()
}

func (v *GroupView) foldSelectedGroup() {
	items := v.groupList.Items()
	next := v.selectedIndex + 1
	for i := next; i < len(items); i++ {
		if items[next] == nil {
			break
		}
		node := items[next].(*host.Node[host.TreeNode])
		if node.Depth > v.selectedGroup.Depth {
			v.groupList.RemoveItem(next)
		}
	}
}
