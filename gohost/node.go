package gohost

import "github.com/charmbracelet/bubbles/list"

type Node interface {
	list.Item
	GetID() string
	GetParentID() string
}

type TreeNode[T Node] struct {
	Node     T
	Children []*TreeNode[T]
	Depth    int
	IsFolded bool
}

func NewTreeNode[T Node](data T, depth int) *TreeNode[T] {
	return &TreeNode[T]{
		Node:     data,
		Children: make([]*TreeNode[T], 0),
		Depth:    depth,
		IsFolded: true,
	}
}

func (n *TreeNode[T]) FilterValue() string {
	return n.Node.FilterValue()
}

func (n *TreeNode[T]) GetID() string {
	return n.Node.GetID()
}

func (n *TreeNode[T]) GetParentID() string {
	return n.Node.GetParentID()
}
