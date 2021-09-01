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
	Hosts   map[string]*Host
	Groups  map[string][]*Host
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
		Hosts:   map[string]*Host{},
		Groups:  map[string][]*Host{},
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
		node := NewHost(file.Name(), path.Join(m.HostDir, file.Name()))
		m.addHost(node)
		m.addGroup(node)
	}
	return nil
}

func (m *manager) PrintGroups() {
	if len(m.Groups) == 0 {
		fmt.Println("no host group")
		return
	}
	fmt.Println("host groups:")
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

func (m *manager) PrintHosts() {
	if len(Manager.Hosts) == 0 {
		fmt.Println("no host file")
	}
	fmt.Println("host files:")
	for name, node := range Manager.Hosts {
		fmt.Printf("  %s (%s)\n", name, node.FileName)
	}
}

func (m *manager) CreateNewHost(name string, groups []string) {
	if _, exist := m.Hosts[name]; exist {
		display.Err(fmt.Errorf("host file '%s' is existed\n", name))
		return
	}
	filePath := m.hostPath(m.genHostName(name, groups))
	err := editor.Open(filePath)
	if err != nil {
		fmt.Printf("failed to create file '%s'\n", filePath)
	}
}

func (m *manager) ChangeHostName(name string, newName string) {
	h, exist := m.Hosts[name]
	if !exist {
		display.Err(fmt.Errorf("host file '%s' is not existed\n", name))
		return
	}
	newHostName := m.genHostName(newName, h.Groups)
	if err := os.Rename(h.Path, m.hostPath(newHostName)); err != nil {
		display.Err(err)
	}
	fmt.Printf("renamed '%s' to '%s'\n", h.Name, newName)
}

func (m *manager) ChangeGroups(name string, newGroups []string) {
	h, exist := m.Hosts[name]
	if !exist {
		display.Err(fmt.Errorf("host file '%s' is not existed\n", name))
		return
	}
	newFile := m.genHostName(name, newGroups)
	if err := os.Rename(h.Path, m.hostPath(newFile)); err != nil {
		display.Err(err)
	}
	fmt.Printf("chanaged group '%v' to '%v\n", h.Groups, newGroups)
}

func (m *manager) EditHostFile(name string) error {
	node, exist := m.Hosts[name]
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

func (m *manager) addHost(host *Host) {
	m.Hosts[host.Name] = host
}

func (m *manager) addGroup(host *Host) {
	for _, group := range host.Groups {
		m.Groups[group] = append(m.Groups[group], host)
	}
}

func (m *manager) printNodes() {
	for name, node := range m.Hosts {
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

func (m *manager) genHostName(name string, groups []string) string {
	var sb strings.Builder
	if len(groups) > 0 {
		sb.WriteString(strings.Join(groups, SpGroup))
		sb.WriteString(SpGroup)
	}
	sb.WriteString(name)
	return sb.String()
}

func (m *manager) hostPath(fileName string) string {
	return path.Join(m.HostDir, fileName)
}
