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
	SpGroup = "_"
	HostDir = ".gohost"
)

type Nodes = map[string]*Node
type Groups = map[string][]*Node

type Switcher struct {
	HomeDir string
	HostDir string
}

var h *Switcher

func init() {
	curr, err := user.Current()
	if err != nil {
		panic(err)
	}
	h = &Switcher{
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

func (h *Switcher) LoadHostNodes() (Nodes, Groups) {
	allNodes := Nodes{}
	allGroups := Groups{}
	err := filepath.Walk(h.HostDir, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		node := NewNode(info.Name(), filePath)
		name, groups := node.NameGroups()
		allNodes[name] = node
		for _, group := range groups {
			allGroups[group] = append(allGroups[group], node)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	printNodes(allNodes)
	printGroups(allGroups)
	return allNodes, allGroups
}

func printNodes(nodes Nodes) {
	for name, node := range nodes {
		fmt.Printf("name %s, node %+v\n", name, node)
	}
}

func printGroups(groups Groups) {
	for group, nodes := range groups {
		fmt.Println("group", group)
		for _, node := range nodes {
			fmt.Printf("node %+v\n", node)
		}
	}
}
