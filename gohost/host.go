package gohost

import (
	"github.com/timshannon/bolthold"
	"gohost/util"
)

type Host interface {
	Node
	GetContent() []byte
	SetContent([]byte)
	IsEnabled() bool
	IsEditable() bool
}

func (s *Service) SaveHost(host Host) error {
	if err := s.store.Insert(host.GetID(), host); err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateHost(host Host) error {
	if err := s.store.Update(host.GetID(), host); err != nil {
		return err
	}
	return nil
}

func (s *Service) LoadHosts(groupID string) []Host {
	return s.loadLocalHosts(groupID)
}

func (s *Service) LoadHostNodes(groupID string) []*TreeNode {
	groupNode := s.nodes[groupID]
	if groupNode == nil {
		return nil
	}
	hostNodeDepth := groupNode.Depth + 1
	hosts := s.LoadHosts(groupID)
	hostNodes := make([]*TreeNode, 0, len(hosts))
	for _, host := range hosts {
		node := NewTreeNode(host)
		node.SetDepth(hostNodeDepth)
		hostNodes = append(hostNodes, node)
	}
	return hostNodes
}

func (s *Service) loadLocalHosts(groupID string) []Host {
	var hosts []*LocalHost
	if err := s.store.FindNullable(&hosts, bolthold.Where("GroupID").Eq(groupID)); err != nil {
		panic(err)
	}
	return util.WrapSlice[Host](hosts)
}
