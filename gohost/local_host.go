package gohost

import (
	"github.com/timshannon/bolthold"
	"gohost/db"
	"gohost/util"
)

var _ Host = (*LocalHost)(nil)

type LocalHost struct {
	ID      db.ID `boltholdKey:"ID"`
	GroupID db.ID
	Name    string
	Content []byte
	Desc    string
	Flag    int
}

func (h *LocalHost) SetFlag(flag int) {
	h.Flag = flag
}

func (h *LocalHost) GetFlag() int {
	return h.Flag
}

func (h *LocalHost) Title() string {
	return h.Name
}

func (h *LocalHost) Description() string {
	return h.Desc
}

func (h *LocalHost) IsEditable() bool {
	return true
}

func (h *LocalHost) GetID() db.ID {
	return h.ID
}

func (h *LocalHost) GetContent() []byte {
	return h.Content
}

func (h *LocalHost) SetContent(content []byte) {
	h.Content = content
}

func (h *LocalHost) FilterValue() string {
	return h.Name
}

func (h *LocalHost) GetParentID() db.ID {
	return h.GroupID
}

func (s *Service) loadLocalHosts(groupID db.ID) []Host {
	var hosts []*LocalHost
	if err := s.store.FindNullable(&hosts, bolthold.Where("GroupID").Eq(groupID)); err != nil {
		panic(err)
	}
	return util.WrapSlice[Host](hosts)
}

func (s *Service) loadLocalHostNodesByParent(parent *TreeNode) []*TreeNode {
	var localHosts []*LocalHost
	if err := s.store.FindNullable(&localHosts, bolthold.Where("GroupID").Eq(parent.GetID())); err != nil {
		panic(err)
	}
	nodes := make([]*TreeNode, len(localHosts))
	for i := range localHosts {
		node := NewTreeNode(localHosts[i])
		node.SetParent(parent)
		nodes[i] = node
	}
	return nodes
}
func (s *Service) DeleteLocalHost(id db.ID) error {
	return s.store.Delete(id, &LocalHost{})
}
