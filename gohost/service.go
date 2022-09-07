package gohost

import (
	"gohost/config"
	"gohost/db"
	"os"
	"sort"
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
		nodes:       make(map[string]*TreeNode[Node], 0),
		tree:        make([]*TreeNode[Node], 0),
		SysHostNode: NewTreeNode[Node](SysHost(), 0),
	}
}

type Service struct {
	store       *db.Store
	nodes       map[string]*TreeNode[Node]
	tree        []*TreeNode[Node]
	SysHostNode *TreeNode[Node]
}

func (s *Service) Tree() []*TreeNode[Node] {
	return s.tree
}

func (s *Service) cacheNodes(nodes []*TreeNode[Node]) {
	for _, node := range nodes {
		s.nodes[node.GetID()] = node
	}
}

func (s *Service) buildTree() {
	// Build tree
	for _, node := range s.nodes {
		p, exist := s.nodes[node.Node.GetParentID()]
		if !exist {
			s.tree = append(s.tree, node)
			continue
		}
		node.Depth = p.Depth + 1
		p.Children = append(p.Children, node)
	}
	// Bfs to set depth
	sort.Slice(s.tree, func(i, j int) bool {
		return s.tree[i].Node.GetID() < s.tree[j].GetID()
	})
	nodes := s.tree
	depth := 0
	for len(nodes) > 0 {
		for _, node := range nodes {
			node.Depth = depth
			sort.Slice(node.Children, func(i, j int) bool {
				return node.Children[i].GetID() < node.Children[j].GetID()
			})
			nodes = append(nodes, node.Children...)
			nodes = nodes[1:]
		}
		depth++
	}
}

func (s *Service) Load() {
	nodes := []*TreeNode[Node]{s.SysHostNode}
	nodes = append(nodes, s.loadGroupNodes()...)
	s.cacheNodes(nodes)
	s.buildTree()
}

func (s *Service) ChildNodes(nodeID string) []*TreeNode[Node] {
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
