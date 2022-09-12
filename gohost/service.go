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
		service.LoadRootNodes()
	})
	return service
}

func NewService() *Service {
	s := &Service{
		store: db.Instance(),
		nodes: make(map[db.ID]*TreeNode, 0),
	}
	s.tree = &TreeNode{
		parent:   nil,
		children: make([]*TreeNode, 0),
		depth:    -1,
		isFolded: false,
	}
	sysHostNode := NewTreeNode(SysHostInstance())
	sysHostNode.SetParent(s.tree)
	s.SysHostNode = sysHostNode
	s.nodes[0] = s.tree
	s.nodes[1] = s.SysHostNode
	return s
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
	nodes = append(nodes, s.LoadNodesByParentID(0)...)
	s.cacheNodes(nodes)
	s.buildTree(nodes)
}

func (s *Service) TreeNodesAsItem() []list.Item {
	nodes := make([]list.Item, 0)
	s.treeNodesAsItem([]*TreeNode{s.tree}, &nodes)
	return nodes
}

func (s *Service) treeNodesAsItem(nodes []*TreeNode, res *[]list.Item) {
	if len(nodes) == 0 {
		return
	}
	for _, node := range nodes {
		if node != s.tree {
			*res = append(*res, node)
		}
		if !node.IsFolded() {
			s.treeNodesAsItem(node.Children(), res)
		}
	}
}

func (s *Service) LoadRootNodes() []*TreeNode {
	return s.LoadNodesByParentID(0)
}

func (s *Service) LoadNodesByParentID(parentID db.ID) []*TreeNode {
	parent := s.nodes[parentID]
	if parent == nil {
		panic("Cache the parent node first before load nodes by parent id")
	}
	var nodes []*TreeNode
	if parentID == 0 {
		nodes = append(nodes, s.SysHostNode)
	}
	nodes = append(nodes, s.loadGroupNodesByParent(parent)...)
	nodes = append(nodes, s.loadLocalHostNodesByParent(parent)...)
	parent.SetChildren(nodes)
	s.cacheNodes(nodes)
	return nodes
}

func (s *Service) RemoveNodesByParentID(parentID db.ID) {
	node := s.nodes[parentID]
	if node == nil {
		panic("node is not cached when trying to remove nodes by parent id")
	}
	node.SetChildren(nil)
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
