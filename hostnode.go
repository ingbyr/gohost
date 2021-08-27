/*
 @Author: ingbyr
*/

package main

import "fmt"

type HostNode struct {
	Name  string
	IsDir bool
	Sub   []*HostNode
}

func (hn *HostNode) AddSub(other *HostNode) {
	hn.Sub = append(hn.Sub, other)
}

func (hn *HostNode) Print() {
	hn.print("")

}

func (hn *HostNode) print(prefix string) {
	fmt.Println(prefix, hn.Name, hn.IsDir)
	for _, sub := range hn.Sub {
		sub.print("    " + prefix)
	}
}
