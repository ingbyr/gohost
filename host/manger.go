/*
 @Author: ingbyr
*/

package host

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ingbyr/gohost/conf"
	"github.com/ingbyr/gohost/display"
	"github.com/ingbyr/gohost/editor"
	"github.com/ingbyr/gohost/util"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type manager struct {
	BaseDir      string
	BaseHostFile string
	Hosts        map[string]*Host
	Groups       map[string][]*Host
}

var Manager *manager

func init() {
	Manager = &manager{
		BaseDir:      conf.BaseDir,
		BaseHostFile: conf.BaseHostFile,
		Hosts:        map[string]*Host{},
		Groups:       map[string][]*Host{},
	}
	if _, err := os.Stat(Manager.BaseDir); os.IsNotExist(err) {
		if err := os.Mkdir(Manager.BaseDir, os.ModePerm); err != nil {
			panic("can not create dir " + Manager.BaseDir)
		}
		fmt.Println("create host dir", Manager.BaseDir)
	}
	if _, err := os.Stat(Manager.BaseHostFile); os.IsNotExist(err) {
		if _, err := os.Create(Manager.BaseHostFile); err != nil {
			panic("can not create base host file")
		}
		fmt.Println("create base host file", Manager.BaseHostFile)
	}
	Manager.LoadHosts()
}

func (m *manager) LoadHosts() {
	files, err := ioutil.ReadDir(m.BaseDir)
	if err != nil {
		display.ErrExit(fmt.Errorf("failed to load gohost dir"))
	}
	for _, file := range files {
		// skip dir and .* files
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}
		node := NewHost(file.Name(), path.Join(m.BaseDir, file.Name()))
		m.addHost(node)
		m.addGroup(node)
	}
}

func (m *manager) PrintGroups() {
	if len(m.Groups) == 0 {
		fmt.Println("no host group")
		return
	}

	header := []string{"Group", "Hosts"}
	data := make([][]string, 0, len(m.Groups))
	for group, hosts := range m.Groups {
		var hsb strings.Builder
		for _, host := range hosts {
			hsb.WriteString(host.Name)
			hsb.WriteString(", ")
		}
		data = append(data, []string{group, hsb.String()[:hsb.Len()-2]})
	}
	display.Table(header, data)
}

func (m *manager) PrintHosts() {
	if len(Manager.Hosts) == 0 {
		fmt.Println("no host file")
		return
	}

	header := []string{"Host", "Groups"}
	data := make([][]string, 0, len(m.Groups))
	for name, node := range Manager.Hosts {
		data = append(data, []string{name, node.GroupsAsStr()})
	}
	display.Table(header, data)
}

func (m *manager) PrintGroup(hostName string) {
	host := m.mustHost(hostName)
	header := []string{"Host", "Groups"}
	data := [][]string{
		{hostName, host.GroupsAsStr()},
	}
	display.Table(header, data)
}

func (m *manager) DeleteGroups(delGroups []string) {
	deleted := make([]string, 0)
	for _, delGroup := range delGroups {
		if hosts, exist := m.Groups[delGroup]; exist {
			// delete hosts which belongs to delGroup
			for _, host := range hosts {
				_ = os.Remove(host.Path)
			}
			deleted = append(deleted, delGroup)
		}
	}
	fmt.Printf("deleted group %s\n", strings.Join(deleted, ","))
}

func (m *manager) DeleteHostGroups(hostName string, delGroups []string) {
	host := m.mustHost(hostName)
	newGroups, removedGroups := util.SliceSub(host.Groups, delGroups)
	hostName = m.hostName(hostName, newGroups)
	err := os.Rename(host.Path, m.fullFilePath(hostName))
	if err != nil {
		display.ErrExit(fmt.Errorf("failed to delete groups"))
	}
	host.Groups = m.groupsFromHostName(hostName)
	fmt.Printf("removed groups '%s'\n", strings.Join(removedGroups, ", "))
}

func (m *manager) AddGroup(hostName string, groups []string) {
	host := m.mustHost(hostName)
	newGroups, addGroups := util.SliceUnion(host.Groups, groups)
	err := os.Rename(host.Path, m.fullFilePath(m.hostName(hostName, newGroups)))
	if err != nil {
		display.ErrExit(fmt.Errorf("failed to delete groups"))
	}
	host.Groups = newGroups
	fmt.Printf("added groups '%s'\n", strings.Join(addGroups, ", "))
}

func (m *manager) CreateNewHost(name string, groups []string) {
	if _, exist := m.Hosts[name]; exist {
		display.ErrExit(fmt.Errorf("host file '%s' is existed\n", name))
		return
	}
	filePath := m.fullFilePath(m.hostName(name, groups))
	err := editor.OpenByVim(filePath)
	if err != nil {
		fmt.Printf("failed to create file '%s'\n", filePath)
	}
}

func (m *manager) DeleteHostsByNames(hostNames []string) {
	deleted := make([]string, 0)
	for _, hostName := range hostNames {
		if host, exist := m.Hosts[hostName]; exist {
			err := os.Remove(host.Path)
			if err != nil {
				display.ErrExit(err)
				continue
			}
			deleted = append(deleted, host.Name)
		}
	}
	fmt.Printf("deleted host %s\n", strings.Join(deleted, ","))
}

func (m *manager) ChangeHostName(hostName string, newHostName string) {
	h := m.mustHost(hostName)
	_newHostName := m.hostName(newHostName, h.Groups)
	if err := os.Rename(h.Path, m.fullFilePath(_newHostName)); err != nil {
		display.ErrExit(err)
	}
	fmt.Printf("renamed '%s' to '%s'\n", h.Name, newHostName)
}

func (m *manager) EditHostFile(hostName string) {
	host := m.mustHost(hostName)
	if err := editor.OpenByVim(host.Path); err != nil {
		display.ErrExit(err)
	}
}

func (m *manager) ApplyGroup(group string) {
	hosts, exist := m.Groups[group]
	if !exist {
		display.ErrExit(fmt.Errorf("not found group '%s'", group))
		return
	}
	combinedHostContent := m.combineHosts(hosts, "# Auto generated from group "+group)
	combinedHost := m.fullFilePath(conf.TmpCombinedHost)
	if err := ioutil.WriteFile(combinedHost, combinedHostContent, 0664); err != nil {
		display.ErrExit(err)
	}
	if err := os.Rename(combinedHost, sysHost); err != nil {
		display.ErrExit(err)
	}
	fmt.Printf("applied group '%s' to system host:\n", group)
	m.PrintSysHost(10)
}

func (m *manager) PrintSysHost(max int) {
	host, err := os.Open(sysHost)
	if err != nil {
		panic(err)
	}
	defer host.Close()
	scanner := bufio.NewScanner(host)
	curr := 0
	for scanner.Scan() {
		if max > 0 && curr == max {
			break
		}
		curr++
		fmt.Println(scanner.Text())
	}
	if scanner.Scan() {
		fmt.Println("...")
	}
}

func (m *manager) host(hostName string) (*Host, bool) {
	host, exist := m.Hosts[hostName]
	if !exist {
		display.ErrExit(fmt.Errorf("host file '%s' is not existed\n", hostName))
		return nil, exist
	}
	return host, exist
}

func (m *manager) mustHost(hostName string) *Host {
	host, exist := m.Hosts[hostName]
	if !exist {
		display.ErrExit(fmt.Errorf("host file '%s' is not existed\n", hostName))
		os.Exit(0)
	}
	return host
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

func (m *manager) hostName(name string, groups []string) string {
	if groups == nil || len(groups) == 0 {
		// use same name as group if no groups specified
		groups = append(groups, name)
	}
	return strings.Join(append(groups, name), conf.SepGroupInFile)
}

func (m *manager) groupsFromHostName(hostName string) []string {
	groups := strings.Split(hostName, conf.SepGroupInFile)
	return groups[:len(groups)-1]
}

func (m *manager) fullFilePath(fileName string) string {
	return path.Join(m.BaseDir, fileName)
}

func (m *manager) combineHosts(hosts []*Host, head string) []byte {
	var b bytes.Buffer
	b.WriteString(head)
	b.WriteByte('\n')
	b.WriteByte('\n')
	for _, host := range hosts {
		file, err := os.Open(host.Path)
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(file)
		b.WriteString("# Host section from " + host.Name + "\n")
		for scanner.Scan() {
			b.Write(scanner.Bytes())
			b.WriteByte('\n')
		}
		b.WriteByte('\n')
		_ = file.Close()
	}
	return b.Bytes()
}
