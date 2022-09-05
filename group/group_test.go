package group

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gohost/config"
	"testing"
)

func TestGroupService_Save(t *testing.T) {
	defer config.store.Close()
	a := assert.New(t)
	gs := NewService()
	groups := []Group{
		{
			ID:     1,
			Parent: 0,
			Name:   "g1",
			Desc:   "desc1",
		},
		{
			ID:     2,
			Parent: 0,
			Name:   "g2",
			Desc:   "desc2",
		},
		{
			ID:     3,
			Parent: 1,
			Name:   "g13",
			Desc:   "desc3",
		},
		{
			ID:     4,
			Parent: 3,
			Name:   "g34",
			Desc:   "desc4",
		},
	}
	for _, g := range groups {
		if err := gs.Save(g); err != nil {
			a.NoError(err)
		}
	}
	savedGroups, err := gs.LoadGroups()
	a.NoError(err)
	for i := range savedGroups {
		fmt.Println("wtf", savedGroups[i])
	}
}

func TestGroupService_BuildTree(t *testing.T) {
	a := assert.New(t)
	gs := NewService()
	groups, err := gs.LoadGroups()
	a.NoError(err)
	for i := range groups {
		fmt.Println(groups[i])
	}
	gs.BuildTree(groups)
	for _, node := range gs.Tree {
		fmt.Println(node)
	}
}
