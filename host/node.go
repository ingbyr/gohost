package host

import "github.com/charmbracelet/bubbles/list"

type Type uint

const (
	GroupType Type = iota
	LocalHostType
	RemoteHostType
)

type TreeNode interface {
	list.Item
	GetID() uint
	GetParentID() uint
}

type Node[T TreeNode] struct {
	Data     T
	Children []*Node[T]
	Depth    int
	Type     Type
	IsFold   bool
}

func NewNode[T TreeNode](data T, depth int) *Node[T] {
	return &Node[T]{
		Data:     data,
		Children: make([]*Node[T], 0),
		Depth:    depth,
		IsFold:   true,
	}
}

func (n *Node[T]) FilterValue() string {
	return n.FilterValue()
}

func (n *Node[T]) GetID() uint {
	return n.Data.GetID()
}

func (n *Node[T]) GetParentID() uint {
	return n.Data.GetParentID()
}
