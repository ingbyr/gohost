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

type service struct {
	Groups map[uint]*Node
	Tree   []*Node
}

func NewService() *service {
	return &service{
		Groups: make(map[uint]*Node, 0),
		Tree:   make([]*Node, 0),
	}
}

func (gs *service) LoadGroups() ([]Group, error) {
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

func (gs *service) BuildTree(groups []Group) {
	if len(gs.Groups) == 0 {
		return
	}
	for _, group := range groups {
		gs.Groups[group.ID] = NewGroupNode(group)
	}
	for _, node := range gs.Groups {
		p, exist := gs.Groups[node.Parent]
		if !exist {
			gs.Tree = append(gs.Tree, node)
			continue
		}
		p.Children = append(p.Children, node)
	}
}

func (gs *service) Save(group Group) error {
	if _, exist := gs.Groups[group.ID]; exist {
		return ErrGroupExist
	}
	err := store.Store().Insert(group.ID, group)
	if err != nil {
		return err
	}
	gs.Groups[group.ID] = NewGroupNode(group)
	return nil
}
