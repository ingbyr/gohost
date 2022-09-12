package gohost

type Host interface {
	Node
	GetContent() []byte
	SetContent([]byte)
	IsEnabled() bool
	IsEditable() bool
}

func (s *Service) SaveHost(host Host) error {
	if err := s.store.Insert(s.extractID(host), host); err != nil {
		return err
	}
	return nil
}

func (s *Service) SaveHostNode(hostNode *TreeNode) error {
	host := hostNode.Node.(Host)
	if err := s.SaveHost(host); err != nil {
		return err
	}
	s.nodes[hostNode.GetID()] = hostNode
	return nil
}

func (s *Service) UpdateHost(host Host) error {
	if err := s.store.Update(host.GetID(), host); err != nil {
		return err
	}
	return nil
}

