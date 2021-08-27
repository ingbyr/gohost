/*
 @Author: ingbyr
*/

package main

import "fmt"

type HostNode struct {
	Name  string
	Path  string
	IsDir bool
	Sub   []*HostNode
}

func NewHostNode(name string, path string, isDir bool) *HostNode {
	return &HostNode{
		Name:  name,
		Path:  path,
		IsDir: isDir,
		Sub:   make([]*HostNode, 0),
	}
}

func (hn *HostNode) AddSub(other *HostNode) {
	hn.Sub = append(hn.Sub, other)
}

func (hn *HostNode) Print() {
	hn.print("")

}

func (hn *HostNode) print(prefix string) {
	fmt.Println(prefix, hn.Name)
	for _, sub := range hn.Sub {
		sub.print("    " + prefix)
	}
}
