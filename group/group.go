package group

import (
	"errors"
)

var (
	ErrGroupExist = errors.New("group is already existed")
)

type Group struct {
	ID     uint `boltholdKey:"ID"`
	Parent uint
	Name   string
	Desc   string
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
