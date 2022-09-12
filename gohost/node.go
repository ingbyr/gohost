package gohost

import (
	"github.com/charmbracelet/bubbles/list"
	"gohost/db"
)

type Node interface {
	list.DefaultItem
	GetID() db.ID
	GetParentID() db.ID
}

type TreeNode struct {
	Node
	parent   *TreeNode
	children []*TreeNode
	depth    int
	isFolded bool
}

func NewTreeNode(node Node) *TreeNode {
	return &TreeNode{
		Node:     node,
		parent:   nil,
		children: make([]*TreeNode, 0),
		depth:    0,
		isFolded: true,
	}
}

func (n *TreeNode) GetID() db.ID {
	if n.Node == nil {
		return 0
	}
	return n.Node.GetID()
}

func (n *TreeNode) Parent() *TreeNode {
	return n.parent
}

func (n *TreeNode) SetParent(parent *TreeNode) {
	n.parent = parent
	n.depth = parent.depth + 1
	parent.children = append(parent.children, n)
}

func (n *TreeNode) Children() []*TreeNode {
	return n.children
}

func (n *TreeNode) SetChildren(children []*TreeNode) {
	n.children = children
}

func (n *TreeNode) RemoveChild(child *TreeNode) {
	for i := range n.children {
		if n.children[i] == child {
			n.SetChildren(append(n.children[:i], n.children[i+1:]...))
			return
		}
	}
}

func (n *TreeNode) Depth() int {
	return n.depth
}

func (n *TreeNode) SetDepth(depth int) {
	n.depth = depth
}

func (n *TreeNode) IsFolded() bool {
	return n.isFolded
}

func (n *TreeNode) SetFolded(isFolded bool) {
	n.isFolded = isFolded
}
