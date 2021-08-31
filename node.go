/*
 @Author: ingbyr
*/

package main

import "strings"

type Node struct {
	Name string
	Path string
}

func NewNode(name string, path string) *Node {
	return &Node{
		Name: name,
		Path: path,
	}
}

func (n *Node) NameGroups() (name string, groups []string) {
	groups = strings.Split(n.Name, SpGroup)
	if len(groups) < 1 {
		panic("not valid file name " + n.Path)
	}
	name = groups[len(groups)-1]
	return
}
