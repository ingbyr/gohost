package gohost

import (
	"github.com/timshannon/bolthold"
	"gohost/util"
)

type Host interface {
	TreeNode
	GetID() uint
	GetName() string
	GetContent() []byte
	GetDesc() string
	GetGroupID() uint
}

func (s *Service) SaveHost(host Host) error {
	if err := s.store.Insert(host.GetDesc(), host); err != nil {
		return err
	}
	return nil
}

func (s *Service) LoadHosts(groupID uint) []Host {
	return s.loadLocalHosts(groupID)
}

func (s *Service) LoadHostNodes(groupID uint) []*Node[TreeNode] {
	groupNode := s.nodes[groupID]
	if groupNode == nil {
		return nil
	}
	hostNodeDepth := groupNode.Depth + 1
	hosts := s.LoadHosts(groupID)
	hostNodes := make([]*Node[TreeNode], 0, len(hosts))
	for _, host := range hosts {
		node := NewNode[TreeNode](host, hostNodeDepth)
		hostNodes = append(hostNodes, node)
	}
	return hostNodes
}

func (s *Service) loadLocalHosts(groupID uint) []Host {
	var hosts []*LocalHost
	if err := s.store.FindNullable(&hosts, bolthold.Where("GroupID").Eq(groupID)); err != nil {
		panic(err)
	}
	return util.WrapSlice[Host](hosts)
}
