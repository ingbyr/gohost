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
	groups := strings.Split(fileName, SpGroup)
	name := groups[len(groups)-1]
	return &Node{
		Name:   name,
		FileName:  fileName,
		Path:   path,
		Groups: groups,
	}
}
