/*
 @Author: ingbyr
*/

package host

import (
	"fmt"
	"github.com/ingbyr/gohost/editor"
	"io/fs"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

const (
	SpGroup      = "_"
	Dir          = ".gohost"
	DefaultGroup = "default"
)

type manager struct {
	HomeDir string
	HostDir string
	Nodes   map[string]*Node
	Groups  map[string][]*Node
}

var Manager *manager

func init() {
	curr, err := user.Current()
	if err != nil {
		panic(err)
	}
	Manager = &manager{
		HomeDir: curr.HomeDir,
		HostDir: path.Join(curr.HomeDir, Dir),
		Nodes:   map[string]*Node{},
		Groups:  map[string][]*Node{},
	}
	if _, err := os.Stat(Manager.HostDir); os.IsNotExist(err) {
		if err := os.Mkdir(Manager.HostDir, os.ModePerm); err != nil {
			panic("can not create dir " + Manager.HostDir)
		}
		fmt.Println("create host dir at", Manager.HostDir)
	}
	if err := Manager.LoadHostNodes(); err != nil {
		panic(err)
	}
}

func (m *manager) LoadHostNodes() error {
	return filepath.Walk(m.HostDir, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		node := NewNode(info.Name(), filePath)
		m.AddNode(node)
		m.AddGroupNode(node)
		return nil
	})
}

func (m *manager) AddNode(node *Node) {
	m.Nodes[node.Name] = node
}

func (m *manager) AddGroupNode(node *Node) {
	for _, group := range node.Groups {
		m.Groups[group] = append(m.Groups[group], node)
	}
}

func (m *manager) PrintGroups() {
	fmt.Println("All groups:")
	for group, nodes := range Manager.Groups {
		var sb strings.Builder
		for _, node := range nodes {
			sb.WriteString(node.Name)
			sb.WriteString(", ")
		}
		nodeNames := sb.String()[:sb.Len()-2]
		fmt.Printf("  %s (%s)\n", group, nodeNames)
	}
}

func (m *manager) PrintHostNodes() {
	fmt.Println("All hosts:")
	for name, node := range Manager.Nodes {
		fmt.Printf("  %s (%s)\n", name, node.FileName)
	}
}

func (m *manager) EditHostFile(name string) error {
	node, exist := m.Nodes[name]
	if !exist {
		return fmt.Errorf("not found host file '%s'", name)
	}
	editor.Open(node.Path)
	return nil
}

func (m *manager) printNodes() {
	for name, node := range m.Nodes {
		fmt.Printf("name %s, node %+v\n", name, node)
	}
}

func (m *manager) printGroups() {
	for group, nodes := range m.Groups {
		fmt.Println("group", group)
		for _, node := range nodes {
			fmt.Printf("node %+v\n", node)
		}
	}
}
