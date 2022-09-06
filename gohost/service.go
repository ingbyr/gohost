package gohost

import (
	"gohost/db"
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
		store: db.Instance(),
		nodes: make(map[uint]*Node[TreeNode], 0),
		tree:  make([]*Node[TreeNode], 0),
	}
}

type Service struct {
	store *db.Store
	nodes map[uint]*Node[TreeNode]
	tree  []*Node[TreeNode]
}

func (s *Service) Tree() []*Node[TreeNode] {
	return s.tree
}

func (s *Service) cacheNodes(nodes []*Node[TreeNode]) {
	for _, node := range nodes {
		s.nodes[node.GetID()] = node
	}
}

func (s *Service) buildTree() {
	// Build tree
	for _, node := range s.nodes {
		p, exist := s.nodes[node.Data.GetParentID()]
		if !exist {
			s.tree = append(s.tree, node)
			continue
		}
		node.Depth = p.Depth + 1
		p.Children = append(p.Children, node)
	}
	// Bfs to set depth
	sort.Slice(s.tree, func(i, j int) bool {
		return s.tree[i].Data.GetID() < s.tree[j].GetID()
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
	nodes := make([]*Node[TreeNode], 0)
	nodes = append(nodes, s.loadGroupNodes()...)
	s.cacheNodes(nodes)
	s.buildTree()
}

func (s *Service) ChildNodes(nodeID uint) []*Node[TreeNode] {
	return s.nodes[nodeID].Children
}
