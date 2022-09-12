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
}

func newNodeItemDelegate() *nodeItemDelegate {
	return &nodeItemDelegate{
		selectedStyle: styles.FocusedFormItem,
		normalStyle:   styles.UnfocusedFormItem,
	}
}

func (d *nodeItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	node, ok := item.(*gohost.TreeNode)
	if !ok {
		return
	}
	var str string
	switch node.Node.(type) {
	case *gohost.Group:
		var icon string
		if node.IsFolded() {
			icon = "📁"
		} else {
			icon = "📂"
		}
		str = strings.Repeat(" ", node.Depth()) + icon + node.Title()
	case *gohost.SysHost:
		str = strings.Repeat(" ", node.Depth()) + "🐝" + node.Title()
	case *gohost.LocalHost:
		str = strings.Repeat(" ", node.Depth()) + "📑" + node.Title()
	case *gohost.RemoteHost:
		str = strings.Repeat(" ", node.Depth()) + "🌐" + node.Title()
	}
	if m.Index() == index {
		str = d.selectedStyle.Render("> " + str)
	} else {
		str = d.normalStyle.Render("  " + str)
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

// TreeView is tui helpView for nodes tree
type TreeView struct {
	model         *Model
	nodeList      list.Model
	width, height int
}

func NewTreeView(model *Model) *TreeView {
	// Create nodes list helpView
	nodes := svc.TreeNodeItem()
	nodeList := list.New(nodes, newNodeItemDelegate(), 0, 0)
	nodeList.Title = "gohost"
	nodeList.SetShowStatusBar(false)
	nodeList.SetShowHelp(false)
	nodeList.Select(0)

	return &TreeView{
		model:    model,
		nodeList: nodeList,
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
				selectedNode := v.SelectedNode()
				switch node := selectedNode.Node.(type) {
				case *gohost.Group:
					selectedNode.FlipFolded()
					cmd = v.RefreshTreeNodes()
				case gohost.Host:
					v.onHostNodeSelected(node, &cmds)
				}
			case key.Matches(m, keys.New):
				v.model.switchState(nodeViewState)
			}
		} else {
			// Disable key
			msg = nil
		}
	}
	v.nodeList, cmd = v.nodeList.Update(msg)
	log.Debug(fmt.Sprintf("cursor at %d, selected item %v",
		v.nodeList.Cursor(), v.nodeList.SelectedItem().FilterValue()))
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

func (v *TreeView) SelectedNode() *gohost.TreeNode {
	return v.nodeList.SelectedItem().(*gohost.TreeNode)
}

func (v *TreeView) onHostNodeSelected(host gohost.Host, cmds *[]tea.Cmd) {
	log.Debug("select host: " + host.Title())
	v.model.switchState(editorViewState)
	v.model.editorView.SetHost(host)
}
