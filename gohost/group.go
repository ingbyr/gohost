package gohost

import (
	"github.com/timshannon/bolthold"
	"gohost/db"
)

var (
	//ErrGroupExist      = errors.New("group is already existed")
	_ Node = (*Group)(nil)
)

type Group struct {
	ID       db.ID `boltholdKey:"ID"`
	ParentID db.ID
	Name     string
	Desc     string
}

func (g *Group) Title() string {
	return "[G] " + g.Name
}
func (g *Group) Description() string {
	return g.Desc
}
func (g *Group) FilterValue() string {
	return g.Name
}

func (g *Group) GetID() db.ID {
	return g.ID
}

func (g *Group) GetParentID() db.ID {
	return g.ParentID
}

func (s *Service) loadGroups() []*Group {
	var groups []*Group
	if err := s.store.FindNullable(&groups, &bolthold.Query{}); err != nil {
		panic(err)
	}
	return groups
}

func (s *Service) loadGroupNodes() []*TreeNode {
	groups := s.loadGroups()
	groupNodes := make([]*TreeNode, 0, len(groups))
	for _, group := range groups {
		groupNodes = append(groupNodes, NewTreeNode(group))
	}
	return groupNodes
}

func (s *Service) SaveGroup(group *Group) error {
	err := s.store.Insert(s.extractID(group), group)
	if err != nil {
		return err
	}
	// FIXME set correct depth
	s.nodes[group.ID] = NewTreeNode(group)
	return nil
}
