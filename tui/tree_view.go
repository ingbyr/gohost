package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohost/gohost"
	"gohost/log"
	"gohost/tui/keys"
	"gohost/tui/styles"
	"io"
	"strings"
)

// nodeItemDelegate is item delegate for groupItem
type nodeItemDelegate struct {
	selectedStyle lipgloss.Style
	normalStyle   lipgloss.Style
	width         int
}

func newNodeItemDelegate() *nodeItemDelegate {
	return &nodeItemDelegate{
		selectedStyle: styles.FocusedFormItem,
		normalStyle:   styles.UnfocusedFormItem,
		width:         0,
	}
}

func (d *nodeItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	node, ok := item.(*gohost.TreeNode)
	if !ok {
		return
	}

	var str string
	var enableIcon string
	if node.IsEnabled() {
		enableIcon = " ✔"
	} else {
		enableIcon = ""
	}

	var nodeIcon string
	switch node.Node.(type) {
	case *gohost.Group:
		if node.IsFolded() {
			nodeIcon = "📁 "
		} else {
			nodeIcon = "📂 "
		}
	case *gohost.SysHost:
		nodeIcon = "👻 "
	case *gohost.LocalHost:
		nodeIcon = "📝 "
	case *gohost.RemoteHost:
		nodeIcon = "🌐 "
	}

	var cursor string
	if m.Index() == index {
		cursor = "> "
	} else {
		cursor = "  "
	}
	str = cursor + strings.Repeat(" ", node.Depth()) + nodeIcon + node.Title() + enableIcon

	strLen := lipgloss.Width(str)
	if strLen > d.width {
		if d.width <= 3 {
			str = strings.Repeat(".", d.width)
		} else {
			str = str[:d.width-3] + "..."
		}
	} else {
		str = str + strings.Repeat(" ", d.width-strLen)
	}

	if m.Index() == index {
		str = d.selectedStyle.Render(str)
	} else {
		str = d.normalStyle.Render(str)
	}

	_, _ = fmt.Fprint(w, str)
}

func (d *nodeItemDelegate) Height() int {
	return 1
}

func (d *nodeItemDelegate) Spacing() int {
	return 0
}

func (d *nodeItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d *nodeItemDelegate) SetWidth(width int) {
	d.width = width
}

// TreeView is tui helpView for nodes tree
type TreeView struct {
	model            *Model
	nodeList         list.Model
	nodeItemDelegate *nodeItemDelegate
	width, height    int
}

func NewTreeView(model *Model) *TreeView {
	// Create nodes list helpView
	nodes := svc.TreeNodesAsItem()
	delegate := newNodeItemDelegate()
	nodeList := list.New(nodes, delegate, 0, 0)
	nodeList.Title = "gohost"
	nodeList.SetShowStatusBar(false)
	nodeList.SetShowHelp(false)
	// FIXME height is wrong when show pagination
	nodeList.SetShowPagination(false)
	// TODO Filter mode is not completable yet
	nodeList.SetFilteringEnabled(false)
	nodeList.Select(0)

	return &TreeView{
		model:            model,
		nodeList:         nodeList,
		nodeItemDelegate: delegate,
	}
}

func (v *TreeView) Init() tea.Cmd {
	v.model.setShortHelp(StateTreeView, []key.Binding{keys.Create, keys.Delete, keys.Apply, keys.Save, keys.ForceQuit})
	v.model.setFullHelp(StateTreeView, append(v.nodeList.FullHelp(), []key.Binding{keys.Create}))
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
	case RefreshTreeViewMsg, AppliedNewHostMsg:
		v.RefreshTreeNodes()
	case tea.KeyMsg:
		selectedNode := v.SelectedNode()
		switch {
		case key.Matches(m, keys.Esc):
			return v, nil
		case key.Matches(m, keys.Enter):
			host, isHostNode := selectedNode.Node.(gohost.Host)
			cmd = func() tea.Msg {
				svc.UpdateFoldOfNode(selectedNode, !selectedNode.IsFolded())
				if !isHostNode {
					return RefreshTreeViewMsg{}
				}
				// Switch to editor view if selected host node
				log.Debug("select host: " + host.Title())
				v.model.switchState(StateEditorView)
				v.model.editorView.SetHostNode(selectedNode)
				return RefreshTreeViewMsg{}
			}
		case key.Matches(m, keys.Create):
			cmd = func() tea.Msg {
				v.model.switchState(StateNodeView)
				svc.UpdateFoldOfNode(selectedNode, false)
				return RefreshTreeViewMsg{}
			}
		case key.Matches(m, keys.Delete):
			cmd = func() tea.Msg {
				if err := svc.DeleteNodeRecursively(selectedNode); err != nil {
					panic(err)
				}
				return RefreshTreeViewMsg{}
			}
		case key.Matches(m, keys.Apply):
			cmd = func() tea.Msg {
				svc.UpdateEnabledOfNode(selectedNode, !selectedNode.IsEnabled())
				hostContent := svc.CombineEnabledHosts()
				svc.ApplyHost(hostContent)
				return AppliedNewHostMsg{}
			}
		}
	}
	cmds = append(cmds, cmd)

	v.nodeList, cmd = v.nodeList.Update(msg)
	//log.Debug(fmt.Sprintf("cursor at %d, selected item %v",
	//	v.nodeList.Cursor(), v.nodeList.SelectedItem().FilterValue()))
	cmds = append(cmds, cmd)
	return v, tea.Batch(cmds...)
}

func (v *TreeView) View() string {
	return v.nodeList.View()
}

func (v *TreeView) SetWidth(width int) {
	v.width = width
	v.nodeList.SetWidth(v.width)
	v.nodeItemDelegate.SetWidth(v.width)
}

func (v *TreeView) SetHeight(height int) {
	v.nodeList.SetHeight(height)
	v.height = height
}

func (v *TreeView) RefreshTreeNodes() tea.Cmd {
	return v.nodeList.SetItems(svc.TreeNodesAsItem())
}

func (v *TreeView) SelectedNode() *gohost.TreeNode {
	return v.nodeList.SelectedItem().(*gohost.TreeNode)
}
