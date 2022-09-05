package group

type Node struct {
	*Group
	Children []*Node
	Depth    int
	IsFold   bool
}

func NewGroupNode(group Group, depth int) *Node {
	return &Node{
		Group:    &group,
		Children: make([]*Node, 0),
		Depth:    depth,
		IsFold:   true,
	}
}

func (n *Node) IsParent(o *Node) bool {
	for _, child := range n.Children {
		if child == o {
			return true
		}
	}
	return false
}
