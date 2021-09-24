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
	baseHost    *Host
	hosts       map[string]*Host
	_hosts      []string
	groups      map[string][]*Host
	_groups     []string
	noGroupHost []*Host
	fs          myfs.HostFs
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
	}

	// setup default fs
	M.SetFs(myfs.NewOsFs())
}

func (m *manager) SetFs(newFs myfs.HostFs) {
	m.fs = newFs

	// create base dir
	if _, err := m.fs.Stat(conf.BaseDir); m.fs.IsNotExist(err) {
		if err := m.fs.MkdirAll(conf.BaseDir, myfs.Perm644); err != nil {
			display.Panic("can not create dir "+conf.BaseDir, err)
		}
	}

	// create base host file
	if _, err := m.fs.Stat(m.baseHost.FilePath); m.fs.IsNotExist(err) {
		var content bytes.Buffer
		content.WriteString("127.0.0.1 localhost")
		content.WriteString(NewLine)
		content.WriteString("::1 localhost")
		content.WriteString(NewLine)
		if err := m.fs.WriteFile(m.baseHost.FilePath, content.Bytes(), 0644); err != nil {
			display.Panic("can not create base host file", err)
		}
	}

	m.LoadHosts()
}

func (m *manager) LoadHosts() {
	// reset map
	m.hosts = make(map[string]*Host)
	m._hosts = make([]string, 0)
	m.groups = make(map[string][]*Host)
	m._groups = make([]string, 0)
	m.noGroupHost = make([]*Host, 0)

	files, err := m.fs.ReadDir(conf.BaseDir)
	if err != nil {
		display.ErrExit(fmt.Errorf("failed to load gohost dir"))
	}

	// load host files
	for _, file := range files {
		// skip dir and files started with '.'
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}
		// create host
		host := NewHostByFileName(file.Name())
		// add host
		m.hosts[host.Name] = host
		// add groups
		if host.hasGroups() {
			for _, group := range host.Groups {
				m.groups[group] = append(m.groups[group], host)
			}
		} else {
			m.noGroupHost = append(m.noGroupHost, host)
		}
	}

	// sort hosts and groups
	for hostName := range m.hosts {
		m._hosts = append(m._hosts, hostName)
	}
	for groupName, group := range m.groups {
		m._groups = append(m._groups, groupName)
		sort.Slice(group, func(i, j int) bool {
			return group[i].Name < group[j].Name
		})
	}
	sort.Strings(m._hosts)
	sort.Strings(m._groups)
	sort.Slice(m.noGroupHost, func(i, j int) bool {
		return m.noGroupHost[i].Name < m.noGroupHost[j].Name
	})
}

func (m *manager) HasNoGroupHost() bool {
	return len(m.noGroupHost) > 0
}

func (m *manager) PrintGroup(hostName string) {
	host := m.mustHost(hostName)
	header := []string{"Host", "Groups"}
	data := [][]string{
		{hostName, host.GroupsAsStr()},
	}
	display.Table(header, data)
}

func (m *manager) DisplayGroups() {
	if len(m.groups) == 0 {
		fmt.Println("no host group")
		return
	}
	header := []string{"Group", "Hosts"}
	data := make([][]string, 0, len(m.groups))
	for _, group := range m._groups {
		hosts := m.groups[group]
		data = append(data, []string{group, HostsToStr(hosts)})
	}
	display.Table(header, data)
	if m.HasNoGroupHost() {
		fmt.Printf("no group hosts [%v]\n", HostsToStr(m.noGroupHost))
		fmt.Printf("please add groups before using hosts\n")
	}
}

func (m *manager) DisplayHosts() {
	if len(m.hosts) == 0 {
		fmt.Println("no host file")
		return
	}
	header := []string{"Host", "Groups"}
	data := make([][]string, 0, len(m.groups))
	for _, host := range m._hosts {
		data = append(data, []string{host, m.hosts[host].GroupsAsStr()})
	}
	display.Table(header, data)
}

func (m *manager) DeleteGroups(delGroups []string) {
	deleted := make([]string, 0)
	for _, delGroup := range delGroups {
		if m.DeleteGroup(delGroup) {
			deleted = append(deleted, delGroup)
		}
	}
	fmt.Printf("deleted group [%s]\n", strings.Join(deleted, ","))
}

func (m *manager) DeleteGroup(group string) bool {
	hosts, exist := m.groups[group]
	if !exist {
		return false
	}
	for _, host := range hosts {
		// delete host which has no group left
		if host.RemoveGroup(group) {
			oldFilePath := host.FilePath
			host.GenAutoFields()
			if err := m.fs.Rename(oldFilePath, host.FilePath); err != nil {
				display.ErrExit(err)
			}
		}
	}
	return true
}

func (m *manager) DeleteHostGroups(hostName string, delGroups []string) {
	host := m.mustHost(hostName)
	newGroups, removedGroups := util.SliceSub(host.Groups, delGroups)
	newHost := NewHostByNameGroups(hostName, newGroups)
	err := m.fs.Rename(host.FilePath, newHost.FilePath)
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
	err := m.fs.Rename(host.FilePath, newHost.FilePath)
	if err != nil {
		display.ErrExit(fmt.Errorf("failed to delete groups"))
	}
	m.hosts[newHost.Name] = newHost
	fmt.Printf("added groups [%s]\n", strings.Join(addGroups, ", "))
}

func (m *manager) CreateNewHost(name string, groups []string, edit bool) {
	if name == m.baseHost.Name {
		display.ErrExit(fmt.Errorf("host file '%s' already exists\n", name))
	}
	if _, exist := m.hosts[name]; exist {
		display.ErrExit(fmt.Errorf("host file '%s' already exists\n", name))
	}
	host := NewHostByNameGroups(name, groups)
	if edit {
		if err := editor.OpenByVim(host.FilePath); err != nil {
			display.ErrExit(fmt.Errorf("failed to create file '%s'\n", host.FilePath))
		}
	} else {
		if err := m.fs.WriteFile(host.FilePath, []byte(""), myfs.Perm644); err != nil {
			display.ErrExit(fmt.Errorf("can not create %s file\n", host.FilePath))
		}
	}
}

func (m *manager) DeleteHosts(hostNames []string) {
	deleted := make([]string, 0)
	for _, hostName := range hostNames {
		if host, exist := m.hosts[hostName]; exist {
			err := m.fs.Remove(host.FilePath)
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
	if err := m.fs.Rename(h.FilePath, newHost.FilePath); err != nil {
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
	if err := m.fs.Rename(combinedHost.FilePath, SysHost); err != nil {
		display.ErrExit(err)
	}
	fmt.Printf("applied group '%s' to system host:\n", group)

	// display system host
	m.PrintSysHost(10)
}

func (m *manager) PrintSysHost(max int) {
	host, err := m.fs.Open(SysHost)
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

func (m *manager) group(groupName string) ([]*Host, bool) {
	group, ok := m.groups[groupName]
	if !ok {
		return nil, false
	}
	return group, ok
}

func (m *manager) printHosts() {
	if len(m._hosts) != len(m.hosts) {
		panic("the size of _hosts and hosts is not equal")
	}
	fmt.Printf("All Hosts\n")
	for _, host := range m._hosts {
		fmt.Printf("\t[host] %+v\n", m.hosts[host])
	}
}

func (m *manager) printGroups() {
	if len(m._groups) != len(m.groups) {
		panic("the size of _groups and groups is not equal")
	}
	fmt.Printf("All Groups\n")
	for _, group := range m._groups {
		fmt.Println("\t[group]", group)
		for _, host := range m.groups[group] {
			fmt.Printf("\t\t[host] %+v\n", host)
		}
	}
}

func (m *manager) combineHosts(hosts []*Host, head string) []byte {
	var b bytes.Buffer
	b.WriteString(head)
	b.WriteString(NewLine + NewLine)
	for _, host := range hosts {
		file, err := m.fs.Open(host.FilePath)
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