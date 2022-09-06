package host

import (
	"errors"
)

var (
	ErrGroupExist = errors.New("nodes is already existed")
)

type Group struct {
	ID       uint `boltholdKey:"ID"`
	ParentID uint
	Name     string
	Desc     string
}

func (g Group) Title() string {
	return g.Name
}
func (g Group) Description() string {
	return g.Desc
}
func (g Group) FilterValue() string {
	return g.Name
}

func (g Group) GetID() uint {
	return g.ID
}

func (g Group) GetParentID() uint {
	return g.ParentID
}
