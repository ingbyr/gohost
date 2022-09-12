package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"gohost/gohost"
	"gohost/log"
	"gohost/tui/keys"
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
	spaces := strings.Repeat(" ", node.Depth())
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
	model            *Model
	nodeList         list.Model
	selectedNode     *gohost.TreeNode
	selectedNodeType *NodeType
	selectedIndex    int
	width, height    int
}

func NewTreeView(model *Model) *TreeView {
	// Create nodes list helpView
	//nodeList := list.New(groups, groupItemDelegate{}, 0, 0)
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	nodeList := list.New(svc.TreeNodeItem(), delegate, 0, 0)
	nodeList.Title = "Groups"
	nodeList.SetShowStatusBar(false)
	nodeList.SetShowHelp(false)

	return &TreeView{
		model:            model,
		nodeList:         nodeList,
		selectedNode:     svc.SysHostNode,
		selectedNodeType: NodeSysHost,
	}
}

func (v *TreeView) Init() tea.Cmd {
	v.model.setShortHelp(treeViewState, v.nodeList.ShortHelp())
	v.model.setFullHelp(treeViewState, v.nodeList.FullHelp())
	return nil
}

func (v *TreeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		v.SetWidth(m.Width)
		v.SetHeight(m.Height)
		log.Debug(fmt.Sprintf("tree view w %d h %d", v.nodeList.Width(), v.nodeList.Height()))
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
						v.selectedNodeType = NodeGroup
						v.selectedNode.FlipFolded()
						cmd = v.RefreshTreeNodes()
					case *gohost.SysHost:
						v.selectedNodeType = NodeSysHost
						v.onHostNodeSelected(&cmds)
					case *gohost.LocalHost:
						v.selectedNodeType = NodeLocalHost
						v.onHostNodeSelected(&cmds)
					case *gohost.RemoteHost:
						v.selectedNodeType = NodeRemoteHost
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

func (v *TreeView) SetWidth(width int) {
	v.nodeList.SetWidth(width)
	v.width = width
}

func (v *TreeView) SetHeight(height int) {
	v.nodeList.SetHeight(height)
	v.height = height
}

func (v *TreeView) RefreshTreeNodes() tea.Cmd {
	return v.nodeList.SetItems(svc.TreeNodeItem())
}

func (v *TreeView) onHostNodeSelected(cmds *[]tea.Cmd) {
	selectedHost := v.selectedNode.Node.(gohost.Host)
	log.Debug("select host: " + selectedHost.Title())
	v.model.switchState(editorViewState)
	v.model.editorView.SetHost(selectedHost)
}
