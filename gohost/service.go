package gohost

import (
	"github.com/charmbracelet/bubbles/list"
	"gohost/config"
	"gohost/db"
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

func (s *Service) DeleteNode(node *TreeNode) error {
	if node == nil || node.Node == nil || node == s.tree || node == s.SysHostNode {
		return nil
	}
	if err := s.store.Delete(node.GetID(), node.Node); err != nil {
		return err
	}
	s.nodes[node.GetID()] = nil
	if node.Parent() != nil {
		node.Parent().RemoveChild(node)
	}
	return nil
}

func (s *Service) DeleteNodeRecursively(node *TreeNode) error {
	if err := s.DeleteNode(node); err != nil {
		return err
	}
	for len(node.Children()) > 0 {
		if err := s.DeleteNodeRecursively(node.Children()[0]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) SaveNode(node *TreeNode) error {
	id := node.GetID()
	if id == 0 {
		id = s.store.NextID()
	}
	return s.store.Insert(id, node.Node)
}

func (s *Service) UpdateNode(node *TreeNode) {
	if err := s.store.Update(node.GetID(), node.Node); err != nil {
		panic(err)
	}
}

func (s *Service) UpdateFoldOfNode(treeNode *TreeNode, folded bool) {
	if treeNode == nil || treeNode == s.SysHostNode {
		return
	}
	treeNode.SetFolded(folded)
	if folded {
		s.RemoveNodesByParent(treeNode)
	} else {
		s.LoadNodesByParent(treeNode, true)
	}
	s.UpdateNode(treeNode)
}

func (s *Service) UpdateEnabledOfNode(node *TreeNode, enabled bool) {
	if node == nil || node == s.SysHostNode {
		return
	}
	node.SetEnabled(enabled)
	s.UpdateNode(node)
	s.LoadNodesByParent(node, false)
	for _, child := range node.Children() {
		s.UpdateEnabledOfNode(child, enabled)
	}
}
