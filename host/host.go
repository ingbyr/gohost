/*
 @Author: ingbyr
*/

package host

import (
	"github.com/ingbyr/gohost/conf"
	"github.com/ingbyr/gohost/util"
	"path"
	"strings"
)

type Host struct {
	Name   string
	Groups []string
	// FileName based on Name ang Groups
	FileName string
	// FilePath based on FileName and conf.BaseDir
	FilePath string
}

func (h *Host) GenAutoFields() {
	h.FileName = strings.Join(append(h.Groups, h.Name), conf.SepGroupInFile)
	h.FilePath = path.Join(conf.BaseDir, h.FileName)
}

func (h *Host) RemoveGroup(group string) int {
	h.Groups = util.SliceRemove(h.Groups, group)
	return len(h.Groups)
}

func (h *Host) GroupsAsStr() string {
	return strings.Join(h.Groups, ", ")
}

func NewHostByFileName(fileName string) *Host {
	groups := strings.Split(fileName, conf.SepGroupInFile)
	name := groups[len(groups)-1]
	groups = groups[:len(groups)-1]
	return &Host{
		Name:     name,
		FileName: fileName,
		FilePath: path.Join(conf.BaseDir, fileName),
		Groups:   groups,
	}
}

func NewHostByNameGroups(hostName string, groups []string) *Host {
	// use host name as group if no specified groups
	if len(groups) == 0 {
		groups = append(groups, hostName)
	} else {
		// sort and unique the groups
		groups = util.SortUniqueStringSlice(groups)
	}
	host := &Host{
		Name:   hostName,
		Groups: groups,
	}
	host.GenAutoFields()
	return host
}
