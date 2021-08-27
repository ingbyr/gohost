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
)

const (
	HostDir = ".host"
)

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

func NewHostNode(name string, isDir bool) *HostNode {
	return &HostNode{
		Name:  name,
		IsDir: isDir,
		Sub:   make([]*HostNode, 0),
	}
}

func (h *HostSwitcher) ListAsTree() {
	rootNode := NewHostNode(h.HomeDir, true)
	nodes := map[string]*HostNode{
		rootNode.Name: rootNode,
	}
	err := filepath.Walk(h.HostDir, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		node := NewHostNode(info.Name(), info.IsDir())
		nodes[filePath] = node
		if pNode, exist := nodes[path.Dir(filePath)]; exist {
			pNode.AddSub(node)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	rootNode.Print()
}
