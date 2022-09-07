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
	Children []*TreeNode
	Depth    int
	IsFolded bool

	// depends on Depth
	spaces string
}

func NewTreeNode(node Node) *TreeNode {
	return &TreeNode{
		Node:     node,
		Children: make([]*TreeNode, 0),
		Depth:    0,
		IsFolded: true,
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

func (n *TreeNode) GetParentID() string {
	return n.Node.GetParentID()
}

func (n *TreeNode) SetDepth(depth int) {
	n.Depth = depth
	n.spaces = strings.Repeat(" ", depth)
}
