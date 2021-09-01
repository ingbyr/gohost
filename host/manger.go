/*
 @Author: ingbyr
*/

package host

import (
	"bufio"
	"fmt"
	"github.com/ingbyr/gohost/display"
	"github.com/ingbyr/gohost/editor"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
)

const (
	SpGroup = "_"
	Dir     = ".gohost"
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
		fmt.Println("Create host dir", Manager.HostDir)
	}
	if err := Manager.LoadHostNodes(); err != nil {
		panic(err)
	}
}

func (m *manager) LoadHostNodes() error {
	files, err := ioutil.ReadDir(m.HostDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		node := NewNode(file.Name(), path.Join(m.HostDir, file.Name()))
		m.addNode(node)
		m.addGroupNode(node)
	}
	return nil
}

func (m *manager) PrintGroups() {
	if len(m.Groups) == 0 {
		fmt.Println("No host group")
		return
	}
	fmt.Println("Host groups:")
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
	if len(Manager.Nodes) == 0 {
		fmt.Println("No host file")
	}
	fmt.Println("Host files:")
	for name, node := range Manager.Nodes {
		fmt.Printf("  %s (%s)\n", name, node.FileName)
	}
}

func (m *manager) AddHost(name string, groups []string) {
	if _, exist := m.Nodes[name]; exist {
		display.Err(fmt.Errorf("Host file '%s' is existed\n", name))
		return
	}
	var sb strings.Builder
	if len(groups) > 0 {
		sb.WriteString(strings.Join(groups, SpGroup))
		sb.WriteString(SpGroup)
	}
	sb.WriteString(name)
	file := path.Join(m.HostDir, sb.String())
	err := editor.Open(file)
	if err != nil {
		fmt.Printf("Can not create file '%s'\n", file)
	}
}

func (m *manager) EditHostFile(name string) error {
	node, exist := m.Nodes[name]
	if !exist {
		return fmt.Errorf("not found host file '%s'", name)
	}
	return editor.Open(node.Path)
}

func (m *manager) GenerateHost(group string) ([]byte, error) {
	nodes, exist := m.Groups[group]
	if !exist {
		return nil, fmt.Errorf("not found group '%s'", group)
	}
	hostLines := make(map[string]*Line)
	for _, node := range nodes {
		file, err := os.Open(node.Path)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			l := parseLine(scanner.Text())
			if l != nil {
				hostLines[l.Domain] = l
			}
		}
		_ = file.Close()
	}
	return nil, nil
}

func (m *manager) addNode(node *Node) {
	m.Nodes[node.Name] = node
}

func (m *manager) addGroupNode(node *Node) {
	for _, group := range node.Groups {
		m.Groups[group] = append(m.Groups[group], node)
	}
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
