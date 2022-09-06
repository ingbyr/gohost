package gohost

import (
	"github.com/stretchr/testify/assert"
	"gohost/db"
	"testing"
)

func TestGroupService_Save(t *testing.T) {
	store := db.Instance()
	defer store.Close()
	a := assert.New(t)
	service := GetService()
	groups := []Group{
		{ID: 1, ParentID: 0, Name: "g-1", Desc: "desc1"},
		{ID: 11, ParentID: 1, Name: "g-1-1", Desc: "desc11"},
		{ID: 111, ParentID: 11, Name: "g-1-1-1", Desc: "desc111"},
		{ID: 112, ParentID: 11, Name: "g-1-1-2", Desc: "desc112"},
		{ID: 12, ParentID: 1, Name: "g-1-2", Desc: "desc12"},
		{ID: 13, ParentID: 1, Name: "g-1-3", Desc: "desc13"},
		{ID: 131, ParentID: 13, Name: "g-1-3-1", Desc: "desc131"},
		{ID: 132, ParentID: 13, Name: "g-1-3-2", Desc: "desc132"},
		{ID: 2, ParentID: 0, Name: "g-2", Desc: "desc2"},
		{ID: 3, ParentID: 0, Name: "g-3", Desc: "desc3"},
		{ID: 31, ParentID: 3, Name: "g-3-1", Desc: "desc31"},
		{ID: 32, ParentID: 3, Name: "g-3-2", Desc: "desc32"},
		{ID: 4, ParentID: 0, Name: "g-4", Desc: "desc4"},
	}
	for _, g := range groups {
		if err := service.SaveGroup(g); err != nil {
			a.NoError(err)
		}
	}
}
