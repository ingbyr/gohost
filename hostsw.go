/*
 @Author: ingbyr
*/

package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

const (
	Splitter = "@"
	HostDir  = ".host"
)

type HostNodes = map[string]*HostNode
type HostGroupNodes = map[string][]*HostNode

type HostSwitcher struct {
	HomeDir string
	HostDir string
}

var h *HostSwitcher

func init() {
	curr, err := user.Current()
	if err != nil {
		panic(err)
	}
	h = &HostSwitcher{
		HomeDir: curr.HomeDir,
		HostDir: path.Join(curr.HomeDir, HostDir),
	}
	if _, err := os.Stat(h.HostDir); os.IsNotExist(err) {
		if err := os.Mkdir(h.HostDir, os.ModePerm); err != nil {
			panic("can not create dir " + h.HostDir)
		}
		fmt.Println("create host dir at", h.HostDir)
	}
}

func (h *HostSwitcher) LoadHostNodes() (HostNodes, HostGroupNodes) {
	hostNodes := HostNodes{
		// root node
		h.HomeDir: NewHostNode(h.HomeDir, h.HomeDir, true),
	}
	hostGroupNodes := HostGroupNodes{}
	err := filepath.Walk(h.HostDir, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		node := NewHostNode(info.Name(), filePath, info.IsDir())
		hostNodes[filePath] = node
		if pNode, exist := hostNodes[path.Dir(filePath)]; exist {
			pNode.AddSub(node)
		}
		if !node.IsDir {
			_name := strings.Split(node.Name, Splitter)
			if len(_name) != 2 {
				panic("not valid file name " + node.Name)
			}
			alias := _name[0]
			node.Name = _name[1]
			hostGroupNodes[alias] = append(hostGroupNodes[alias], node)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	h.printHostNodes(hostNodes)
	h.printHostGroupNodes(hostGroupNodes)
	return hostNodes, hostGroupNodes
}

func (h *HostSwitcher) printHostNodes(hn HostNodes) {
	for _, subNode := range hn[h.HomeDir].Sub[0].Sub {
		subNode.Print()
	}
}

func (h *HostSwitcher) printHostGroupNodes(hgn HostGroupNodes) {
	for alias, nodes := range hgn {
		fmt.Println("group", alias)
		for _, node := range nodes {
			fmt.Println(node.Name, node.Path)
		}
	}
}
