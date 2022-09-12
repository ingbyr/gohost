package tui

var (
	NodeGroup      = &NodeType{"Group", "Group contains some hosts"}
	NodeSysHost    = &NodeType{"System Host", "Current system host"}
	NodeLocalHost  = &NodeType{"Local Host", "Host stored in local database"}
	NodeRemoteHost = &NodeType{"Remote Host", "Host from internet"}
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
