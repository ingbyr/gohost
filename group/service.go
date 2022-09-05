package group

import (
	"errors"
	"github.com/timshannon/bolthold"
	"gohost/store"
	"sync"
)

var (
	instance *service
	once     sync.Once
)

func Service() *service {
	once.Do(func() {
		instance = NewService()
	})
	return instance
}

func NewService() *service {
	return &service{
		groups: make(map[uint]*Node, 0),
		tree:   make([]*Node, 0),
	}
}

type service struct {
	groups map[uint]*Node
	tree   []*Node
}

func (gs *service) Tree() []*Node {
	return gs.tree
}

func (gs *service) loadGroups() ([]Group, error) {
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

func (gs *service) buildTree(groups []Group) {
	if len(groups) == 0 {
		return
	}
	for _, group := range groups {
		gs.groups[group.ID] = NewGroupNode(group)
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

func (gs *service) Load() error {
	groups, err := gs.loadGroups()
	if err != nil {
		return err
	}
	gs.buildTree(groups)
	return nil
}

func (gs *service) Save(group Group) error {
	if _, exist := gs.groups[group.ID]; exist {
		return ErrGroupExist
	}
	err := store.Store().Insert(group.ID, group)
	if err != nil {
		return err
	}
	gs.groups[group.ID] = NewGroupNode(group)
	return nil
}
