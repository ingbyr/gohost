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

func (s *Service) SaveGroup(group *Group) error {
	err := s.store.Insert(s.extractID(group), group)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) SaveGroupNode(groupNode *TreeNode) error {
	group := groupNode.Node.(*Group)
	if err := s.SaveGroup(group); err != nil {
		return err
	}
	s.nodes[groupNode.GetID()] = groupNode
	return nil
}

func (s *Service) DeleteGroup(id db.ID) error {
	return s.store.Delete(id, &Group{})
}

func (s *Service) DeleteGroupNode(groupNode *TreeNode) {
	// Delete from db
	if err := s.DeleteGroup(groupNode.GetID()); err != nil {
		panic(err)
	}
	// Delete from parent
	groupNode.Parent().RemoveChild(groupNode)
	// Delete all children
	children := s.LoadNodesByParent(groupNode)
	for i := range children {
		s.DeleteGroupNode(children[i])
	}
	// Delete node cache
	s.nodes[groupNode.GetID()] = nil
}
