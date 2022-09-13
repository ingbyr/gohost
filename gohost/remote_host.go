package gohost

import (
	"gohost/db"
)

var _ Host = (*RemoteHost)(nil)

type RemoteHost struct {
}

func (r *RemoteHost) SetFlag(flag int) {
	//TODO implement me
	panic("implement me")
}

func (r *RemoteHost) GetFlag() int {
	//TODO implement me
	panic("implement me")
}

func (r *RemoteHost) GetContent() []byte {
	//TODO implement me
	panic("implement me")
}

func (r *RemoteHost) SetContent(bytes []byte) {
	//TODO implement me
	panic("implement me")
}

func (r *RemoteHost) IsEditable() bool {
	//TODO implement me
	panic("implement me")
}

func (r *RemoteHost) FilterValue() string {
	//TODO implement me
	panic("implement me")
}

func (r *RemoteHost) Title() string {
	//TODO implement me
	panic("implement me")
}

func (r *RemoteHost) Description() string {
	//TODO implement me
	panic("implement me")
}

func (r *RemoteHost) GetID() db.ID {
	//TODO implement me
	panic("implement me")
}

func (r *RemoteHost) GetParentID() db.ID {
	//TODO implement me
	panic("implement me")
}
