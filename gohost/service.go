package gohost

import (
	"github.com/charmbracelet/bubbles/list"
	"gohost/config"
	"gohost/db"
	"os"
	"sync"
)

var (
	service     *Service
	serviceOnce sync.Once
	cfg         = config.Instance()
)

func GetService() *Service {
	serviceOnce.Do(func() {
		service = NewService()
		service.loadTree()
	})
	return service
}

func NewService() *Service {
	return &Service{
		store: db.Instance(),
		nodes: make(map[db.ID]*TreeNode, 0),
		tree: &TreeNode{
			Node:     &LocalHost{
				ID:      0,
				GroupID: 0,
				Name:    "",
				Content: nil,
				Desc:    "",
				Enabled: false,
			},
			parent:   nil,
			children: make([]*TreeNode, 0),
			depth:    -1,
			isFolded: false,
		},
		SysHostNode: NewTreeNode(SysHostInstance()),
	}
}

type Service struct {
	store       *db.Store
	nodes       map[db.ID]*TreeNode
	tree        *TreeNode
	SysHostNode *TreeNode
}

// Tree the system host tree node is always first
func (s *Service) Tree() *TreeNode {
	return s.tree
}

func (s *Service) cacheNodes(nodes []*TreeNode) {
	for _, node := range nodes {
		s.nodes[node.GetID()] = node
	}
}

func (s *Service) buildTree(nodes []*TreeNode) {
	// Build tree
	for _, node := range nodes {
		p := s.nodes[node.Node.GetParentID()]
		node.SetParent(p)
		p.children = append(p.children, node)
	}
	// Bfs to set depth
	queue := s.tree.children
	depth := 0
	for len(queue) > 0 {
		for _, treeNode := range queue {
			treeNode.SetDepth(depth)
			queue = append(queue, treeNode.children...)
			queue = queue[1:]
		}
		depth++
	}
}

func (s *Service) loadTree() {
	s.nodes[0] = s.tree
	nodes := []*TreeNode{s.SysHostNode}
	nodes = append(nodes, s.loadGroupNodes()...)
	s.cacheNodes(nodes)
	s.buildTree(nodes)
}

func (s *Service) TreeNodeItem() []list.Item {
	nodes := make([]list.Item, 0)
	s.treeNodesItem([]*TreeNode{s.tree}, &nodes)
	return nodes
}

func (s *Service) treeNodesItem(nodes []*TreeNode, res *[]list.Item) {
	if len(nodes) == 0 {
		return
	}
	for _, node := range nodes {
		if node != s.tree {
			*res = append(*res, node)
		}
		if !node.IsFolded() {
			s.treeNodesItem(node.Children(), res)
			*res = append(*res, s.LoadHostNodesItem(node.GetID())...)
		}
	}
}

func (s *Service) Node(nodeID db.ID) *TreeNode {
	if nodeID <= 0 {
		return s.tree
	}
	return s.nodes[nodeID]
}

func (s *Service) ChildNodes(nodeID db.ID) []*TreeNode {
	return s.nodes[nodeID].children
}

// ApplyHost TODO apply host to system
func (s *Service) ApplyHost(hosts []byte) {
	// open system host file
	sysHostFile, err := os.Create(cfg.SysHostFile)
	if err != nil {
		panic(err)
	}
	defer sysHostFile.Close()

	// write hosts to system host file
	if _, err = sysHostFile.Write(hosts); err != nil {
		panic(err)
	}
}

func (s *Service) CombineHost(hosts ...[]byte) []byte {
	// TODO combine host
	return nil
}

func (s *Service) extractID(node Node) db.ID {
	if node.GetID() == 0 {
		return s.store.NextID()
	}
	return node.GetID()
}
