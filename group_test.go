package main

import (
	"fmt"
	"testing"
)

func TestGroupService_BuildTree(t *testing.T) {
	//a := assert.New(t)
	gs := NewGroupService()

	groups := []*Group{
		{
			ID:     1,
			Parent: 0,
			Name:   "g1",
		},
		{
			ID:     2,
			Parent: 0,
			Name:   "g2",
		},
		{
			ID:     3,
			Parent: 1,
			Name:   "g13",
		},
		{
			ID:     4,
			Parent: 3,
			Name:   "g34",
		},
	}
	gs.BuildTree(groups)
	for _, node := range gs.Tree {
		fmt.Println(node)
	}
}
