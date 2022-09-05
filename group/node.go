package group

type Node struct {
	*Group
	Children []*Node
	Depth    int
}

func NewGroupNode(group Group, depth int) *Node {
	return &Node{
		Group:    &group,
		Children: make([]*Node, 0),
		Depth:    depth,
	}
}
