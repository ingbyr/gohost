package gohost

import (
	"github.com/charmbracelet/bubbles/list"
	"strings"
)

type Node interface {
	list.DefaultItem
	GetID() string
	GetParentID() string
}

type TreeNode struct {
	Node
	parent   *TreeNode
	children []*TreeNode
	depth    int
	isFolded bool
	spaces   string // depends on depth
}

func NewTreeNode(node Node) *TreeNode {
	return &TreeNode{
		Node:     node,
		parent:   nil,
		children: make([]*TreeNode, 0),
		depth:    0,
		isFolded: true,
		spaces:   "",
	}
}

func (n *TreeNode) Title() string {
	return n.spaces + n.Node.Title()
}

func (n *TreeNode) Description() string {
	return n.spaces + n.Node.Description()
}

func (n *TreeNode) FilterValue() string {
	return n.Node.FilterValue()
}

func (n *TreeNode) GetID() string {
	return n.Node.GetID()
}

func (n *TreeNode) Parent() *TreeNode {
	return n.parent
}

func (n *TreeNode) Children() []*TreeNode {
	return n.children
}

func (n *TreeNode) Depth() int {
	return n.depth
}

func (n *TreeNode) SetDepth(depth int) {
	n.depth = depth
	n.spaces = strings.Repeat(" ", depth)
}

func (n *TreeNode) IsFolded() bool {
	return n.isFolded
}

func (n *TreeNode) SetFolded(isFolded bool) {
	n.isFolded = isFolded
}
