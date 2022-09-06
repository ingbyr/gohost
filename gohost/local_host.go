package gohost

type LocalHost struct {
	ID      uint
	Name    string
	Content []byte
	Desc    string
	GroupID uint
}

// Implement of Host

func (h *LocalHost) GetID() uint {
	return h.ID
}

func (h *LocalHost) GetName() string {
	return h.Name
}

func (h *LocalHost) GetContent() []byte {
	return h.Content
}

func (h *LocalHost) GetDesc() string {
	return h.Desc
}

func (h *LocalHost) GetGroupID() uint {
	return h.GroupID
}

// Implement of TreeNode

func (h *LocalHost) FilterValue() string {
	return h.Name
}

func (h *LocalHost) GetParentID() uint {
	return h.GroupID
}
