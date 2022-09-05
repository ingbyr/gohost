package group

type Node struct {
	*Group
	Children []*Node
}

func NewGroupNode(group Group) *Node {
	return &Node{
		Group:    &group,
		Children: make([]*Node, 0),
	}
}
