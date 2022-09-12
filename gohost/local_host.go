package gohost

import "gohost/db"

var _ Host = (*LocalHost)(nil)

type LocalHost struct {
	ID      db.ID `boltholdKey:"ID"`
	GroupID db.ID
	Name    string
	Content []byte
	Desc    string
	Enabled bool
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

func (h *LocalHost) IsEnabled() bool {
	return h.Enabled
}

func (h *LocalHost) FilterValue() string {
	return h.Name
}

func (h *LocalHost) GetParentID() db.ID {
	return h.GroupID
}
