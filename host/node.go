/*
 @Author: ingbyr
*/

package host

import "strings"

type Node struct {
	Name   string
	FileName string
	Path   string
	Groups []string
}

func NewNode(fileName string, path string) *Node {
	if fileName == "" {
		panic("empty file name " + path)
	}
	groups := strings.Split(fileName, SpGroup)
	name := groups[len(groups)-1]
	groups = groups[:len(groups) - 1]
	if len(groups) == 0 {
		groups = append(groups, DefaultGroup)
	}
	return &Node{
		Name:   name,
		FileName:  fileName,
		Path:   path,
		Groups: groups,
	}
}
