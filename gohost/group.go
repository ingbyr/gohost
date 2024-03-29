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
	Flag     int
}

func (g *Group) SetFlag(flag int) {
	g.Flag = flag
}

func (g *Group) GetFlag() int {
	return g.Flag
}

func (g *Group) Title() string {
	return g.Name
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

func (s *Service) loadGroupsByParentID(parentID db.ID) []*Group {
	var groups []*Group
	if err := s.store.FindNullable(&groups, bolthold.Where("ParentID").Eq(parentID)); err != nil {
		panic(err)
	}
	return groups
}

func (s *Service) loadGroupNodesByParent(parent *TreeNode) []*TreeNode {
	groups := s.loadGroupsByParentID(parent.GetID())
	nodes := make([]*TreeNode, len(groups))
	for i := range groups {
		node := NewTreeNode(groups[i])
		node.SetParent(parent)
		nodes[i] = node
	}
	return nodes
}
