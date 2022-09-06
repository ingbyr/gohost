package host

import (
	"errors"
	"github.com/timshannon/bolthold"
	"gohost/store"
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
		nodes: make(map[uint]*Node[TreeNode], 0),
		tree:  make([]*Node[TreeNode], 0),
	}
}

type Service struct {
	nodes map[uint]*Node[TreeNode]
	tree  []*Node[TreeNode]
}

func (s *Service) Tree() []*Node[TreeNode] {
	return s.tree
}

func (s *Service) loadGroups() ([]Group, error) {
	var groups []Group
	if err := store.Store().Find(&groups, &bolthold.Query{}); err != nil {
		if errors.Is(bolthold.ErrNotFound, err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return groups, nil
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

func (s *Service) Load() error {
	groups, err := s.loadGroups()
	if err != nil {
		return err
	}
	for _, group := range groups {
		s.nodes[group.ID] = NewNode[TreeNode](group, 0)
	}
	s.buildTree()
	return nil

}
func (s *Service) SaveGroup(group Group) error {
	if _, exist := s.nodes[group.ID]; exist {
		return ErrGroupExist
	}
	err := store.Store().Insert(group.ID, group)
	if err != nil {
		return err
	}
	s.nodes[group.ID] = NewNode[TreeNode](&group, 0)
	return nil
}

func (s *Service) Children(nodeID uint) []*Node[TreeNode] {
	return s.nodes[nodeID].Children
}
