package gohost

import "github.com/charmbracelet/bubbles/list"

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
}

func NewTreeNode(data Node, depth int) *TreeNode {
	return &TreeNode{
		Node:     data,
		Children: make([]*TreeNode, 0),
		Depth:    depth,
		IsFolded: true,
	}
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
