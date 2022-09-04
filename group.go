package main

import (
	"errors"
	"github.com/timshannon/bolthold"
)

var (
	ErrGroupExist = errors.New("group is already existed")
)

type Group struct {
	ID     int
	Parent int
	Name   string
}

type GroupNode struct {
	*Group
	Children []*GroupNode
}

func NewGroupNode(group *Group) *GroupNode {
	return &GroupNode{
		Group:    group,
		Children: make([]*GroupNode, 0),
	}
}

type GroupService struct {
	Groups map[int]*GroupNode
	Tree   []*GroupNode
}

func NewGroupService() *GroupService {
	return &GroupService{
		Groups: make(map[int]*GroupNode, 0),
		Tree:   make([]*GroupNode, 0),
	}
}

func (gs *GroupService) LoadGroups() ([]*Group, error) {
	var groups []*Group
	if err := store.Find(&groups, &bolthold.Query{}); err != nil {
		if errors.Is(bolthold.ErrNotFound, err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return groups, nil
}

func (gs *GroupService) BuildTree(groups []*Group) {
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

func (gs *GroupService) Save(group Group) error {
	if _, exist := gs.Groups[group.ID]; exist {
		return ErrGroupExist
	}
	err := store.Insert(bolthold.NextSequence(), group)
	if err != nil {
		return err
	}
	gs.Groups[group.ID] = NewGroupNode(&group)
	return nil
}
