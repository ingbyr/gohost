package gohost

import "github.com/charmbracelet/bubbles/list"

type TreeNode interface {
	list.Item
	GetID() string
	GetParentID() string
}

type Node[T TreeNode] struct {
	Data     T
	Children []*Node[T]
	Depth    int
	IsFolded bool
}

func NewNode[T TreeNode](data T, depth int) *Node[T] {
	return &Node[T]{
		Data:     data,
		Children: make([]*Node[T], 0),
		Depth:    depth,
		IsFolded: true,
	}
}

func (n *Node[T]) FilterValue() string {
	return n.Data.FilterValue()
}

func (n *Node[T]) GetID() string {
	return n.Data.GetID()
}

func (n *Node[T]) GetParentID() string {
	return n.Data.GetParentID()
}
