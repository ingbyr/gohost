package gohost

import (
	"github.com/charmbracelet/bubbles/list"
	"gohost/config"
	"gohost/db"
	"os"
	"sync"
)

var (
	service     *Service
	serviceOnce sync.Once
	cfg         = config.Instance()
)

func GetService() *Service {
	serviceOnce.Do(func() {
		service = NewService()
		service.LoadNodesByParent(service.tree, true)
	})
	return service
}

func NewService() *Service {
	s := &Service{
		store: db.Instance(),
		nodes: make(map[db.ID]*TreeNode, 0),
	}
	s.tree = &TreeNode{
		Node:     &Group{ID: 0},
		parent:   nil,
		children: make([]*TreeNode, 0),
		depth:    -1,
	}
	sysHostNode := NewTreeNode(SysHostInstance())
	sysHostNode.SetParent(s.tree)
	s.SysHostNode = sysHostNode
	s.nodes[0] = s.tree
	s.nodes[1] = s.SysHostNode
	return s
}

type Service struct {
	store       *db.Store
	nodes       map[db.ID]*TreeNode
	tree        *TreeNode
	SysHostNode *TreeNode
}

// Tree the system host tree node is always first
func (s *Service) Tree() *TreeNode {
	return s.tree
}

func (s *Service) cacheNodes(nodes []*TreeNode) {
	for _, node := range nodes {
		s.nodes[node.GetID()] = node
	}
}

func (s *Service) buildTree(nodes []*TreeNode) {
	// Build tree
	for _, node := range nodes {
		p := s.nodes[node.Node.GetParentID()]
		node.SetParent(p)
	}
	// Bfs to set depth
	queue := s.tree.children
	depth := 0
	for len(queue) > 0 {
		for _, treeNode := range queue {
			treeNode.SetDepth(depth)
			queue = append(queue, treeNode.children...)
			queue = queue[1:]
		}
		depth++
	}
}

func (s *Service) TreeNodesAsItem() []list.Item {
	nodes := make([]list.Item, 0)
	s.treeNodesAsItem([]*TreeNode{s.tree}, &nodes)
	return nodes
}

func (s *Service) treeNodesAsItem(nodes []*TreeNode, res *[]list.Item) {
	if len(nodes) == 0 {
		return
	}
	for _, node := range nodes {
		if node != s.tree {
			*res = append(*res, node)
		}
		if !node.IsFolded() {
			s.treeNodesAsItem(node.children, res)
		}
	}
}

func (s *Service) LoadNodesByParent(parent *TreeNode, considerFold bool) []*TreeNode {
	var children []*TreeNode
	if parent == s.tree {
		children = append(children, s.SysHostNode)
	}
	children = append(children, s.loadGroupNodesByParent(parent)...)
	children = append(children, s.loadLocalHostNodesByParent(parent)...)
	parent.SetChildren(children)
	s.cacheNodes(children)
	nodes := append([]*TreeNode{}, children...)
	for _, node := range children {
		if considerFold && node.IsFolded() {
			continue
		}
		nodes = append(nodes, s.LoadNodesByParent(node, considerFold)...)
	}
	return nodes
}

func (s *Service) RemoveNodesByParent(parent *TreeNode) {
	parent.SetChildren(nil)
}

func (s *Service) extractID(node Node) db.ID {
	if node.GetID() == 0 {
		return s.store.NextID()
	}
	return node.GetID()
}

func (s *Service) DeleteNode(node *TreeNode) {
	if node == nil || node.Node == nil || node == s.SysHostNode {
		return
	}
	if node == s.tree {
		panic("Can not delete dummy root node")
	}
	node.Parent().RemoveChild(node)
	s.nodes[node.GetID()] = nil
	switch node.Node.(type) {
	case *Group:
		if err := s.DeleteGroup(node.GetID()); err != nil {
			panic(err)
		}
		for _, childNode := range node.Children() {
			s.DeleteNode(childNode)
		}
	case *LocalHost:
		if err := s.DeleteLocalHost(node.GetID()); err != nil {
			panic(err)
		}
	case *RemoteHost:

	}
}

func (s *Service) UpdateNode(node *TreeNode) {
	if err := s.store.Update(node.GetID(), node.Node); err != nil {
		panic(err)
	}
}

func (s *Service) ApplyHost(hosts []byte) {
	// Truncate system host file
	sysHostFile, err := os.Create(cfg.SysHostFile)
	if err != nil {
		panic(err)
	}
	defer sysHostFile.Close()
	// Write hosts to system host file
	if _, err = sysHostFile.Write(hosts); err != nil {
		panic(err)
	}
}

func (s *Service) EnableHost() {
	// TODO enable as group node
	// TODO enable as localhost node
}

func (s *Service) UnfoldNode(treeNode *TreeNode) {
	if !s.isFoldableNode(treeNode) {
		return
	}
	treeNode.SetFolded(false)
	s.LoadNodesByParent(treeNode, true)
	s.UpdateGroupNode(treeNode)
}

func (s *Service) FoldNode(treeNode *TreeNode) {
	if !s.isFoldableNode(treeNode) {
		return
	}
	treeNode.SetFolded(true)
	s.RemoveNodesByParent(treeNode)
	s.UpdateGroupNode(treeNode)
}

func (s *Service) isFoldableNode(treeNode *TreeNode) bool {
	if treeNode == nil || treeNode == s.SysHostNode {
		return false
	}
	_, ok := treeNode.Node.(*Group)
	return ok
}

func (s *Service) EnableNode(node *TreeNode) {
	if node == nil || node == s.SysHostNode {
		return
	}
	node.SetEnabled(true)
	s.UpdateNode(node)
	if s.isFoldableNode(node) {
		s.LoadNodesByParent(node, false)
		for _, child := range node.Children() {
			s.EnableNode(child)
		}
	}
}
