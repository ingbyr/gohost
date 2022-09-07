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
)

func GetService() *Service {
	serviceOnce.Do(func() {
		service = NewService()
	})
	return service
}

func NewService() *Service {
	return &Service{
		store:       db.Instance(),
		nodes:       make(map[string]*TreeNode, 0),
		tree:        make([]*TreeNode, 0),
		SysHostNode: NewTreeNode(SysHost(), 0),
	}
}

type Service struct {
	store       *db.Store
	nodes       map[string]*TreeNode
	tree        []*TreeNode
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
		node.Depth = p.Depth + 1
		p.Children = append(p.Children, node)
	}
	// Bfs to set depth
	queue := s.tree
	depth := 0
	for len(queue) > 0 {
		for _, treeNode := range queue {
			treeNode.Depth = depth
			queue = append(queue, treeNode.Children...)
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

func (s *Service) ChildNodes(nodeID string) []*TreeNode {
	return s.nodes[nodeID].Children
}

// ApplyHost TODO apply host to system
func (s *Service) ApplyHost(hosts []byte) {
	// open system host file
	sysHostFile, err := os.Create(config.Instance().SysHostFile)
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
