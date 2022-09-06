package gohost

type LocalHost struct {
	ID      string `boltholdKey:"ID"`
	Name    string
	Content []byte
	Desc    string
	GroupID string
}

// Implement of Host
var _ Host = (*LocalHost)(nil)

func (h *LocalHost) GetID() string {
	return h.ID
}

func (h *LocalHost) GetName() string {
	return h.Name
}

func (h *LocalHost) GetContent() []byte {
	return h.Content
}

func (h *LocalHost) SetContent(content []byte) {
	h.Content = content
}

func (h *LocalHost) GetDesc() string {
	return h.Desc
}

func (h *LocalHost) GetGroupID() string {
	return h.GroupID
}

// Implement of TreeNode
var _ TreeNode = (*LocalHost)(nil)

func (h *LocalHost) FilterValue() string {
	return h.Name
}

func (h *LocalHost) GetParentID() string {
	return h.GroupID
}
