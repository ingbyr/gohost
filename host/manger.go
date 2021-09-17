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
	"github.com/ingbyr/gohost/myfs"
	"github.com/ingbyr/gohost/util"
	"io/ioutil"
	"sort"
	"strings"
)

type manager struct {
	baseHost *Host
	hosts    map[string]*Host
	groups map[string][]*Host
	fs     myfs.HostFs
}

var M *manager

func init() {
	// init manager
	M = &manager{
		baseHost: &Host{
			Name:     conf.BaseHostFileName,
			FileName: conf.BaseHostFileName,
			FilePath: conf.BaseHostFile,
			Groups:   nil,
		},
		hosts:  map[string]*Host{},
		groups: map[string][]*Host{},
		fs:     myfs.NewOsFs(),
	}

	// create base dir
	if _, err := M.fs.Stat(conf.BaseDir); M.fs.IsNotExist(err) {
		if err := M.fs.MkdirAll(conf.BaseDir, myfs.Perm664); err != nil {
			display.Panic("can not create dir "+conf.BaseDir, err)
		}
	}

	// create base host file
	if _, err := M.fs.Stat(M.baseHost.FilePath); M.fs.IsNotExist(err) {
		var content bytes.Buffer
		content.WriteString("127.0.0.1 localhost")
		content.WriteString(NewLine)
		content.WriteString("::1 localhost")
		content.WriteString(NewLine)
		if err := M.fs.WriteFile(M.baseHost.FilePath, content.Bytes(), 0644); err != nil {
			display.Panic("can not create base host file", err)
		}
	}
	M.LoadHosts()
}

func (m *manager) LoadHosts() {
	files, err := M.fs.ReadDir(conf.BaseDir)
	if err != nil {
		display.ErrExit(fmt.Errorf("failed to load gohost dir"))
	}
	// load host files
	for _, file := range files {
		// skip dir and files started with '.'
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}
		host := NewHostByFileName(file.Name())
		// add host
		m.hosts[host.Name] = host
		// add groups
		for _, group := range host.Groups {
			m.groups[group] = append(m.groups[group], host)
		}
	}
}

func (m *manager) PrintGroup(hostName string) {
	host := m.mustHost(hostName)
	header := []string{"Host", "groups"}
	data := [][]string{
		{hostName, host.GroupsAsStr()},
	}
	display.Table(header, data)
}

func (m *manager) PrintGroups() {
	if len(m.groups) == 0 {
		fmt.Println("no host group")
		return
	}
	header := []string{"Group", "hosts"}
	groupNames := make([]string, 0, len(m.groups))
	for groupName := range m.groups {
		groupNames = append(groupNames, groupName)
	}
	sort.Strings(groupNames)
	data := make([][]string, 0, len(m.groups))
	for _, groupName := range groupNames {
		hosts := m.groups[groupName]
		var hsb strings.Builder
		for _, host := range hosts {
			hsb.WriteString(host.Name)
			hsb.WriteString(", ")
		}
		data = append(data, []string{groupName, hsb.String()[:hsb.Len()-2]})
	}
	display.Table(header, data)
}

func (m *manager) PrintHosts() {
	if len(M.hosts) == 0 {
		fmt.Println("no host file")
		return
	}
	header := []string{"Host", "groups"}
	// sort host names
	data := make([][]string, 0, len(m.groups))
	hostNames := make([]string, 0, len(M.hosts))
	for hostName := range M.hosts {
		hostNames = append(hostNames, hostName)
	}
	sort.Strings(hostNames)
	// append display data
	for _, hostName := range hostNames {
		data = append(data, []string{hostName, M.hosts[hostName].GroupsAsStr()})
	}
	display.Table(header, data)
}

func (m *manager) DeleteGroups(delGroups []string) {
	deleted := make([]string, 0)
	for _, delGroup := range delGroups {
		if hosts, exist := m.groups[delGroup]; exist {
			// delete hosts which belongs to delGroup
			for _, host := range hosts {
				_ = M.fs.Remove(host.FilePath)
			}
			deleted = append(deleted, delGroup)
		}
	}
	fmt.Printf("deleted group [%s]\n", strings.Join(deleted, ","))
}

func (m *manager) DeleteHostGroups(hostName string, delGroups []string) {
	host := m.mustHost(hostName)
	newGroups, removedGroups := util.SliceSub(host.Groups, delGroups)
	newHost := NewHostByNameGroups(hostName, newGroups)
	err := M.fs.Rename(host.FilePath, newHost.FilePath)
	if err != nil {
		display.ErrExit(fmt.Errorf("failed to delete groups"))
	}
	m.hosts[newHost.Name] = newHost
	fmt.Printf("removed groups [%s]\n", strings.Join(removedGroups, ", "))
}

func (m *manager) AddGroup(hostName string, groups []string) {
	host := m.mustHost(hostName)
	newGroups, addGroups := util.SliceUnion(host.Groups, groups)
	newHost := NewHostByNameGroups(hostName, newGroups)
	err := M.fs.Rename(host.FilePath, newHost.FilePath)
	if err != nil {
		display.ErrExit(fmt.Errorf("failed to delete groups"))
	}
	m.hosts[newHost.Name] = newHost
	fmt.Printf("added groups [%s]\n", strings.Join(addGroups, ", "))
}

func (m *manager) CreateNewHost(name string, groups []string) {
	if name == m.baseHost.Name {
		display.ErrExit(fmt.Errorf("host file '%s' already exists\n", name))
	}
	if _, exist := m.hosts[name]; exist {
		display.ErrExit(fmt.Errorf("host file '%s' already exists\n", name))
	}
	host := NewHostByNameGroups(name, groups)
	err := editor.OpenByVim(host.FilePath)
	if err != nil {
		fmt.Printf("failed to create file '%s'\n", host.FilePath)
	}
}

func (m *manager) DeleteHosts(hostNames []string) {
	deleted := make([]string, 0)
	for _, hostName := range hostNames {
		if host, exist := m.hosts[hostName]; exist {
			err := M.fs.Remove(host.FilePath)
			if err != nil {
				display.ErrExit(err)
				continue
			}
			deleted = append(deleted, host.Name)
		}
	}
	fmt.Printf("deleted host [%s]\n", strings.Join(deleted, ","))
}

func (m *manager) ChangeHostName(hostName string, newHostName string) {
	if hostName == m.baseHost.Name || newHostName == m.baseHost.Name {
		display.ErrExit(fmt.Errorf("can not change base host file name"))
	}
	h := m.mustHost(hostName)
	newHost := NewHostByNameGroups(newHostName, h.Groups)
	if err := M.fs.Rename(h.FilePath, newHost.FilePath); err != nil {
		display.ErrExit(err)
	}
	fmt.Printf("renamed '%s' to '%s'\n", h.Name, newHostName)
}

func (m *manager) EditHostFile(hostName string) {
	host := m.mustHost(hostName)
	if err := editor.OpenByVim(host.FilePath); err != nil {
		display.ErrExit(err)
	}
}

func (m *manager) ApplyGroup(group string, simulate bool) {
	hosts, exist := m.groups[group]
	if !exist {
		display.ErrExit(fmt.Errorf("not found group '%s'", group))
		return
	}
	hosts = append(hosts, m.baseHost)
	combinedHostContent := m.combineHosts(hosts, "# Auto generated from group "+group)

	// just print
	if simulate {
		fmt.Println(string(combinedHostContent))
		return
	}

	// write to tmp file
	combinedHost := NewHostByNameGroups(conf.TmpCombinedHost, nil)
	if err := ioutil.WriteFile(combinedHost.FilePath, combinedHostContent, 0664); err != nil {
		display.ErrExit(err)
	}

	// replace system host
	if err := M.fs.Rename(combinedHost.FilePath, SysHost); err != nil {
		display.ErrExit(err)
	}
	fmt.Printf("applied group '%s' to system host:\n", group)

	// display system host
	m.PrintSysHost(10)
}

func (m *manager) PrintSysHost(max int) {
	host, err := M.fs.Open(SysHost)
	if err != nil {
		display.Panic("can not read system host file", err)
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
	if hostName == m.baseHost.Name {
		return m.baseHost, true
	}
	host, exist := m.hosts[hostName]
	if !exist {
		display.ErrExit(fmt.Errorf("host file '%s' is not existed\n", hostName))
		return nil, exist
	}
	return host, exist
}

func (m *manager) mustHost(hostName string) *Host {
	if hostName == m.baseHost.Name {
		return m.baseHost
	}
	host, exist := m.hosts[hostName]
	if !exist {
		display.ErrExit(fmt.Errorf("host file '%s' is not existed\n", hostName))
	}
	return host
}

func (m *manager) printNodes() {
	for name, node := range m.hosts {
		fmt.Printf("name %s, node %+v\n", name, node)
	}
}

func (m *manager) printGroups() {
	for group, nodes := range m.groups {
		fmt.Println("group", group)
		for _, node := range nodes {
			fmt.Printf("node %+v\n", node)
		}
	}
}

func (m *manager) combineHosts(hosts []*Host, head string) []byte {
	var b bytes.Buffer
	b.WriteString(head)
	b.WriteString(NewLine + NewLine)
	for _, host := range hosts {
		file, err := M.fs.Open(host.FilePath)
		if err != nil {
			display.Panic("can not combine host", err)
		}
		scanner := bufio.NewScanner(file)
		b.WriteString("# Host section from " + host.Name + NewLine)
		for scanner.Scan() {
			b.Write(scanner.Bytes())
			b.WriteString(NewLine)
		}
		b.WriteString(NewLine)
		_ = file.Close()
	}
	return b.Bytes()
}
