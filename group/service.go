package group

import (
	"errors"
	"github.com/timshannon/bolthold"
	"gohost/store"
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
		groups: make(map[uint]*Node, 0),
		tree:   make([]*Node, 0),
	}
}

type Service struct {
	groups map[uint]*Node
	tree   []*Node
}

func (gs *Service) Tree() []*Node {
	return gs.tree
}

func (gs *Service) Nodes() []*Node {
	nodes := make([]*Node, 0, len(gs.groups))
	for _, node := range gs.tree {
		gs.dfsNode(&nodes, node, 0)
	}
	return nodes
}

func (gs *Service) dfsNode(nodes *[]*Node, node *Node, depth int) {
	node.Depth = depth
	*nodes = append(*nodes, node)
	for _, child := range node.Children {
		gs.dfsNode(nodes, child, depth+1)
	}
}

func (gs *Service) loadGroups() ([]Group, error) {
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

func (gs *Service) buildTree(groups []Group) {
	if len(groups) == 0 {
		return
	}
	for _, group := range groups {
		gs.groups[group.ID] = NewGroupNode(group, 0)
	}
	for _, node := range gs.groups {
		p, exist := gs.groups[node.Parent]
		if !exist {
			gs.tree = append(gs.tree, node)
			continue
		}
		p.Children = append(p.Children, node)
	}
}

func (gs *Service) Load() error {
	groups, err := gs.loadGroups()
	if err != nil {
		return err
	}
	gs.buildTree(groups)
	return nil
}

// Save TODO get parent and set depth
func (gs *Service) Save(group Group) error {
	if _, exist := gs.groups[group.ID]; exist {
		return ErrGroupExist
	}
	err := store.Store().Insert(group.ID, group)
	if err != nil {
		return err
	}
	gs.groups[group.ID] = NewGroupNode(group, 0)
	return nil
}
