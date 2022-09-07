package gohost

import (
	"errors"
	"github.com/timshannon/bolthold"
)

var (
	ErrGroupExist = errors.New("group is already existed")
)

type Group struct {
	ID       string `boltholdKey:"ID"`
	ParentID string
	Name     string
	Desc     string
}

func (g Group) Title() string {
	return g.Name
}
func (g Group) Description() string {
	return g.Desc
}
func (g Group) FilterValue() string {
	return g.Name
}

func (g Group) GetID() string {
	return g.ID
}

func (g Group) GetParentID() string {
	return g.ParentID
}

func (s *Service) loadGroups() []Group {
	var groups []Group
	if err := s.store.FindNullable(&groups, &bolthold.Query{}); err != nil {
		panic(err)
	}
	return groups
}

func (s *Service) loadGroupNodes() []*TreeNode {
	groups := s.loadGroups()
	groupNodes := make([]*TreeNode, 0, len(groups))
	for _, group := range groups {
		groupNodes = append(groupNodes, NewTreeNode(group, 0))
	}
	return groupNodes
}

func (s *Service) SaveGroup(group Group) error {
	if _, exist := s.nodes[group.ID]; exist {
		return ErrGroupExist
	}
	err := s.store.Insert(group.ID, group)
	if err != nil {
		return err
	}
	// FIXME set correct depth
	s.nodes[group.ID] = NewTreeNode(&group, 0)
	return nil
}
