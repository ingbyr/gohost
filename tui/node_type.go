package tui

var (
	GroupNode  = &NodeType{"Group", "Group contains some hosts"}
	LocalHost  = &NodeType{"Local Host", "Host stored in local database"}
	RemoteHost = &NodeType{"Remote Host", "Host from internet"}
)

type NodeType struct {
	Name string
	Desc string
}

func (n *NodeType) FilterValue() string {
	return n.Name
}

func (n *NodeType) Title() string {
	return n.Name
}

func (n *NodeType) Description() string {
	return n.Desc
}
