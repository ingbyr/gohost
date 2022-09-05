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
		{ID: 2, Parent: 0, Name: "g-2", Desc: "desc2"},
		{ID: 3, Parent: 1, Name: "g-1-1", Desc: "desc3"},
		{ID: 4, Parent: 3, Name: "g-1-1-1", Desc: "desc4"},
		{ID: 5, Parent: 3, Name: "g-1-1-2", Desc: "desc5"},
		{ID: 6, Parent: 1, Name: "g-1-2", Desc: "desc6"},
		{ID: 7, Parent: 1, Name: "g-1-3", Desc: "desc7"},
		{ID: 8, Parent: 0, Name: "g-3", Desc: "desc8"},
		{ID: 9, Parent: 0, Name: "g-4", Desc: "desc9"},
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
