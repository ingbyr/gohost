package group

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gohost/store"
	"testing"
)

func TestGroupService_Save(t *testing.T) {
	s := store.Store()
	defer s.Close()
	a := assert.New(t)
	gs := NewService()
	groups := []Group{
		{ID: 1, Parent: 0, Name: "g-1", Desc: "desc1"},
		{ID: 11, Parent: 1, Name: "g-1-1", Desc: "desc11"},
		{ID: 111, Parent: 11, Name: "g-1-1-1", Desc: "desc111"},
		{ID: 112, Parent: 11, Name: "g-1-1-2", Desc: "desc112"},
		{ID: 12, Parent: 1, Name: "g-1-2", Desc: "desc12"},
		{ID: 13, Parent: 1, Name: "g-1-3", Desc: "desc13"},
		{ID: 131, Parent: 13, Name: "g-1-3-1", Desc: "desc131"},
		{ID: 132, Parent: 13, Name: "g-1-3-2", Desc: "desc132"},
		{ID: 2, Parent: 0, Name: "g-2", Desc: "desc2"},
		{ID: 3, Parent: 0, Name: "g-3", Desc: "desc3"},
		{ID: 31, Parent: 3, Name: "g-3-1", Desc: "desc31"},
		{ID: 32, Parent: 3, Name: "g-3-2", Desc: "desc32"},
		{ID: 4, Parent: 0, Name: "g-4", Desc: "desc4"},
	}
	for _, g := range groups {
		if err := gs.Save(g); err != nil {
			a.NoError(err)
		}
	}
	savedGroups, err := gs.loadGroups()
	a.NoError(err)
	for i := range savedGroups {
		fmt.Println("wtf", savedGroups[i])
	}
}

func TestGroupService_BuildTree(t *testing.T) {
	a := assert.New(t)
	gs := NewService()
	groups, err := gs.loadGroups()
	a.NoError(err)
	for i := range groups {
		fmt.Println(groups[i])
	}
	gs.buildTree(groups)
	for _, node := range gs.tree {
		fmt.Println(node)
	}
}
