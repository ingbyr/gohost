package gohost

import (
	"github.com/charmbracelet/bubbles/list"
	"gohost/db"
)

type Node interface {
	list.DefaultItem
	GetID() db.ID
	GetParentID() db.ID
	SetFlag(flag int)
	GetFlag() int
}

const (
	MaskFold = 1 << iota
	MaskEnable
)

type TreeNode struct {
	Node
	parent   *TreeNode
	children []*TreeNode
	depth    int
}

func NewTreeNode(node Node) *TreeNode {
	return &TreeNode{
		Node:     node,
		parent:   nil,
		children: make([]*TreeNode, 0),
		depth:    0,
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
	return n.Node.GetFlag()&MaskFold == MaskFold
}

func (n *TreeNode) SetFolded(folded bool) {
	flag := n.Node.GetFlag()
	if folded {
		n.Node.SetFlag(flag | MaskFold)
	} else {
		n.Node.SetFlag(flag & (^MaskFold))
	}
}

func (n *TreeNode) IsEnabled() bool {
	return n.Node.GetFlag()&MaskEnable == MaskEnable
}

func (n *TreeNode) SetEnabled(enabled bool) {
	flag := n.Node.GetFlag()
	if enabled {
		n.Node.SetFlag(flag | MaskEnable)
	} else {
		n.Node.SetFlag(flag & (^MaskEnable))
	}
}
