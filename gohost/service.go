package gohost

import (
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
		service.Load()
	})
	return service
}

func NewService() *Service {

	return &Service{
		store:       db.Instance(),
		nodes:       make(map[db.ID]*TreeNode, 0),
		tree:        make([]*TreeNode, 0),
		SysHost:     SysHostInstance(),
		SysHostNode: NewTreeNode(SysHostInstance()),
	}
}

type Service struct {
	store       *db.Store
	nodes       map[db.ID]*TreeNode
	tree        []*TreeNode
	SysHost     Host
	SysHostNode *TreeNode
}

// Tree the system host tree node is always first
func (s *Service) Tree() []*TreeNode {
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
		p, exist := s.nodes[node.Node.GetParentID()]
		if !exist {
			s.tree = append(s.tree, node)
			continue
		}
		node.SetDepth(p.depth + 1)
		node.parent = p
		p.children = append(p.children, node)
	}
	// Bfs to set depth
	queue := s.tree
	depth := 0
	for len(queue) > 0 {
		for _, treeNode := range queue {
			treeNode.depth = depth
			queue = append(queue, treeNode.children...)
			queue = queue[1:]
		}
		depth++
	}
}

func (s *Service) Load() {
	nodes := []*TreeNode{s.SysHostNode}
	nodes = append(nodes, s.loadGroupNodes()...)
	s.cacheNodes(nodes)
	s.buildTree(nodes)
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
