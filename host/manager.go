/*
 @Author: ingbyr
*/

package host

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/ingbyr/gohost/config"
	"github.com/ingbyr/gohost/display"
	"github.com/ingbyr/gohost/editor"
	"github.com/ingbyr/gohost/hfs"
	"github.com/ingbyr/gohost/util"
)

type manager struct {
	baseHost    *Host
	hosts       map[string]*Host
	_hosts      []string
	groups      map[string][]*Host
	_groups     []string
	noGroupHost []*Host
	editor      editor.Editor
	config      config.Config
}

var (
	M  *manager
	fs = hfs.H
)

func init() {
	// init config
	_config := config.Config{
		Editor:     editor.Default,
		EditorArgs: editor.DefaultArgs,
	}
	if err := _config.Init(); err != nil {
		display.ErrStrExit("failed to init config file", err)
	}

	// init editor
	_editor := editor.New(_config.Editor, editor.ExtractArgs(_config.EditorArgs))

	// create manager
	M = &manager{
		baseHost: &Host{
			Name:     config.BaseHostFileName,
			FileName: config.BaseHostFileName,
			FilePath: config.BaseHostFile,
			Groups:   nil,
		},
		editor: _editor,
		config: _config,
	}
	M.Init()
}

func (m *manager) Init() {
	// create base host file
	if _, err := fs.Stat(m.baseHost.FilePath); fs.IsNotExist(err) {
		var content bytes.Buffer
		content.WriteString("127.0.0.1 localhost")
		content.WriteString(config.NewLine)
		content.WriteString("::1 localhost")
		content.WriteString(config.NewLine)
		if err := fs.WriteFile(m.baseHost.FilePath, content.Bytes(), 0644); err != nil {
			display.ErrStrExit("can not create base host file", err)
		}
	}
	// load hosts
	m.LoadHosts()
}

func (m *manager) LoadHosts() {
	// reset map
	m.hosts = make(map[string]*Host)
	m._hosts = make([]string, 0)
	m.groups = make(map[string][]*Host)
	m._groups = make([]string, 0)
	m.noGroupHost = make([]*Host, 0)

	files, err := fs.ReadDir(config.BaseDir)
	if err != nil {
		display.ErrStrExit("failed to Init gohost dir")
	}

	// load host files
	for _, file := range files {
		// skip dir and files started with '.'
		if file.IsDir() || !strings.HasSuffix(file.Name(), config.HostFileExt) {
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
			if err := fs.Rename(oldFilePath, host.FilePath); err != nil {
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
	err := fs.Rename(host.FilePath, newHost.FilePath)
	if err != nil {
		display.ErrStrExit("failed to delete groups")
	}
	m.hosts[newHost.Name] = newHost
	fmt.Printf("removed groups [%s]\n", strings.Join(removedGroups, ", "))
}

func (m *manager) AddGroup(hostName string, groups []string) {
	host := m.mustHost(hostName)
	newGroups, addGroups := util.SliceUnion(host.Groups, groups)
	newHost := NewHostByNameGroups(hostName, newGroups)
	err := fs.Rename(host.FilePath, newHost.FilePath)
	if err != nil {
		display.ErrStrExit("failed to delete groups")
	}
	m.hosts[newHost.Name] = newHost
	fmt.Printf("added groups [%s]\n", strings.Join(addGroups, ", "))
}

func (m *manager) CreateNewHost(name string, groups []string, edit bool) {
	if _, exist := m.hosts[name]; exist {
		display.ErrStrExit(fmt.Sprintf("host file '%s' already exists", name))
	}
	host := NewHostByNameGroups(name, groups)
	// create the file before editing
	if err := fs.WriteFile(host.FilePath, []byte(""), hfs.Perm644); err != nil {
		display.ErrStrExit(fmt.Sprintf("failed to create file %s", host.FilePath), err)
	}
	if edit {
		if err := m.editor.Open(host.FilePath); err != nil {
			display.ErrStrExit(fmt.Sprintf("failed to edit file '%s'", host.FilePath), err)
		}
	}
}

func (m *manager) DeleteHosts(hostNames []string) {
	deleted := make([]string, 0)
	for _, hostName := range hostNames {
		if host, exist := m.hosts[hostName]; exist {
			err := fs.Remove(host.FilePath)
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
		display.ErrStrExit("can not change base host file name")
	}
	if _, exist := m.host(newHostName); exist {
		display.ErrStrExit(fmt.Sprintf("host '%s' has been existed", newHostName))
	}
	h := m.mustHost(hostName)
	newHost := NewHostByNameGroups(newHostName, h.Groups)
	if err := fs.Rename(h.FilePath, newHost.FilePath); err != nil {
		display.ErrExit(err)
	}
	fmt.Printf("renamed '%s' to '%s'\n", h.Name, newHostName)
}

func (m *manager) ChangeGroupName(groupName string, newGroupName string) {
	group := m.mustGroup(groupName)
	if groupName == newGroupName {
		return
	}
	for _, host := range group {
		newGroups, _ := util.SliceRemove(host.Groups, groupName)
		newGroups = append(newGroups, newGroupName)
		newGroups = util.SortUniqueStringSlice(newGroups)
		newHost := NewHostByNameGroups(host.Name, newGroups)
		if err := fs.Rename(host.FilePath, newHost.FilePath); err != nil {
			display.ErrExit(err)
		}
	}
	fmt.Printf("rename group '%s' to '%s'\n", groupName, newGroupName)
}

func (m *manager) EditHostFile(hostName string) {
	host := m.mustHost(hostName)
	if err := m.editor.Open(host.FilePath); err != nil {
		display.ErrExit(err)
	}
}

func (m *manager) ApplyGroup(group string, simulate bool) {
	hosts, exist := m.groups[group]
	if !exist {
		display.ErrStrExit(fmt.Sprintf("not found group '%s'", group))
		return
	}
	hosts = append(hosts, m.baseHost)
	combinedHostContent := m.combineHosts(hosts, "# Auto generated from "+group)

	// just print
	if simulate {
		fmt.Println(string(combinedHostContent))
		return
	}

	// open system host file
	sysHost, err := fs.Create(config.SysHost)
	if err != nil {
		display.ErrExit(err)
	}
	defer sysHost.Close()

	// write hosts to system host file
	if _, err = sysHost.Write(combinedHostContent); err != nil {
		display.ErrExit(err)
	}

	// display system host
	m.PrintSysHost(10)
}

func (m *manager) PrintSysHost(max int) {
	host, err := fs.Open(config.SysHost)
	if err != nil {
		display.ErrStrExit("can not read system host file", err)
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

func (m *manager) ChangeConfig(option string, value string) {
	if err := m.config.Change(option, value); err != nil {
		display.ErrStrExit("failed to modify the config file", err)
	}
}

func (m *manager) host(hostName string) (*Host, bool) {
	if hostName == m.baseHost.Name {
		return m.baseHost, true
	}
	host, exist := m.hosts[hostName]
	return host, exist
}

func (m *manager) mustHost(hostName string) *Host {
	host, exist := m.host(hostName)
	if !exist {
		display.ErrStrExit(fmt.Sprintf("host file '%s' is not existed", hostName))
	}
	return host
}

func (m *manager) group(groupName string) ([]*Host, bool) {
	group, exist := m.groups[groupName]
	return group, exist
}

func (m *manager) mustGroup(groupName string) []*Host {
	group, exist := m.group(groupName)
	if !exist {
		display.ErrStrExit(fmt.Sprintf("group '%s' is not existed", groupName))
	}
	return group
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

func (m *manager) combineHosts(hosts []*Host, head string) []byte {
	var b bytes.Buffer
	b.WriteString(head)
	b.WriteString(config.NewLine + config.NewLine)
	for _, host := range hosts {
		file, err := fs.Open(host.FilePath)
		if err != nil {
			display.ErrStrExit("can not combine host", err)
		}
		scanner := bufio.NewScanner(file)
		b.WriteString("# " + host.Name + config.NewLine)
		for scanner.Scan() {
			b.Write(scanner.Bytes())
			b.WriteString(config.NewLine)
		}
		b.WriteString(config.NewLine)
		_ = file.Close()
	}
	return b.Bytes()
}
