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
	node, ok := item.(*gohost.TreeNode)
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
	case *gohost.Group:
		str += fmt.Sprintf("%s[G] %d. %s", spaces, index, node.Name)
	case gohost.Host:
		str += fmt.Sprintf("%s[L] %d. %s", spaces, index, node.Title())
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

// TreeView is tui helpView for nodes tree
type TreeView struct {
	model         *Model
	nodeList      list.Model
	selectedNode  *gohost.TreeNode
	selectedIndex int
	selectedGroup *gohost.Group
	selectedHost  gohost.Host

	service *gohost.Service
}

func NewTreeView(model *Model) *TreeView {
	// Get nodes service
	service := gohost.GetService()
	service.Load()
	treeNodes := service.Tree()
	groups := util.WrapSlice[list.Item](treeNodes)

	// Create nodes list helpView
	//nodeList := list.New(groups, groupItemDelegate{}, 0, 0)
	delegate := list.NewDefaultDelegate()
	groupList := list.New(groups, delegate, 0, 0)
	// TODO add remaining help key
	groupList.Title = "Groups"
	groupList.SetShowHelp(false)

	return &TreeView{
		model:        model,
		nodeList:     groupList,
		selectedNode: treeNodes[0],
		service:      service,
	}
}

func (v *TreeView) Init() tea.Cmd {
	v.model.SetShortHelp(treeViewState, v.nodeList.ShortHelp())
	v.model.SetFullHelp(treeViewState, v.nodeList.FullHelp())
	return nil
}

func (v *TreeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.nodeList.SetHeight(m.Height - v.model.reservedHeight)
		v.nodeList.SetWidth(m.Width / 3)
		v.model.helpView.debug = fmt.Sprintf("w %d h %d, w %d h %d", m.Width, m.Height, v.nodeList.Width(), v.nodeList.Height())
	case tea.KeyMsg:
		if v.model.state == treeViewState {
			switch {
			case key.Matches(m, keys.Enter):
				selectedItem := v.nodeList.SelectedItem()
				if selectedItem != nil {
					v.selectedNode = selectedItem.(*gohost.TreeNode)
					v.selectedIndex = v.nodeList.Index()
					switch v.selectedNode.Node.(type) {
					case *gohost.Group:
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
	v.nodeList, cmd = v.nodeList.Update(msg)
	cmds = append(cmds, cmd)
	return v, tea.Batch(cmds...)
}

func (v *TreeView) View() string {
	return v.nodeList.View()
}

func (v *TreeView) ShortHelp() []key.Binding {
	return v.nodeList.ShortHelp()
}

func (v *TreeView) FullHelp() [][]key.Binding {
	return v.nodeList.FullHelp()
}

func (v *TreeView) onGroupNodeEnterClick(cmds *[]tea.Cmd) {
	v.selectedGroup = v.selectedNode.Node.(*gohost.Group)
	if v.selectedNode.IsFolded {
		v.unfoldSelectedGroup(cmds)
	} else {
		v.foldSelectedGroup()
	}
	v.selectedNode.IsFolded = !v.selectedNode.IsFolded
}

func (v *TreeView) unfoldSelectedGroup(cmds *[]tea.Cmd) {
	subGroups := v.service.ChildNodes(v.selectedNode.GetID())
	idx := v.selectedIndex
	for i := range subGroups {
		idx++
		*cmds = append(*cmds, v.nodeList.InsertItem(idx, subGroups[i]))
	}
	subHosts := v.service.LoadHostNodes(v.selectedNode.GetID())
	for i := range subHosts {
		idx++
		*cmds = append(*cmds, v.nodeList.InsertItem(idx, subHosts[i]))
	}
}

func (v *TreeView) foldSelectedGroup() {
	items := v.nodeList.Items()
	next := v.selectedIndex + 1
	for i := next; i < len(items); i++ {
		if items[next] == nil {
			break
		}
		node := items[next].(*gohost.TreeNode)
		if node.Depth > v.selectedNode.Depth {
			node.IsFolded = true
			v.nodeList.RemoveItem(next)
		} else {
			break
		}
	}
}

func (v *TreeView) onHostNodeSelected(cmds *[]tea.Cmd) {
	v.selectedHost = v.selectedNode.Node.(gohost.Host)
	v.model.log("select host: " + v.selectedHost.Title())
	v.model.switchState(editorViewState)
	v.model.editorView.SetHost(v.selectedHost)
}
