package gohost

var _ Host = (*LocalHost)(nil)

type LocalHost struct {
	ID      string `boltholdKey:"ID"`
	Name    string
	Content []byte
	Desc    string
	GroupID string
	Enabled bool
}

func (h *LocalHost) Title() string {
	return "[L] " + h.Name
}

func (h *LocalHost) Description() string {
	return h.Desc
}

func (h *LocalHost) IsEditable() bool {
	return true
}

func (h *LocalHost) GetID() string {
	return h.ID
}

func (h *LocalHost) GetContent() []byte {
	return h.Content
}

func (h *LocalHost) SetContent(content []byte) {
	h.Content = content
}

func (h *LocalHost) IsEnabled() bool {
	return h.Enabled
}

func (h *LocalHost) FilterValue() string {
	return h.Name
}

func (h *LocalHost) GetParentID() string {
	return h.GroupID
}
