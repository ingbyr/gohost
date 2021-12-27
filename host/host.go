/*
 @Author: ingbyr
*/

package host

import (
	"path/filepath"
	"strings"

	"github.com/ingbyr/gohost/config"
	"github.com/ingbyr/gohost/util"
)

type Host struct {
	Name   string
	Groups []string
	// FileName based on Name ang Groups
	FileName string
	// FilePath based on FileName and config.BaseDir
	FilePath string
}

func (h *Host) GenAutoFields() {
	h.FileName = strings.Join(append(h.Groups, h.Name), config.SepGroupInFile)
	h.FilePath = filepath.Join(config.BaseDir, h.FileName+config.HostFileExt)
}

func (h *Host) RemoveGroup(group string) bool {
	removed := false
	h.Groups, removed = util.SliceRemove(h.Groups, group)
	return removed
}

func (h *Host) GroupsAsStr() string {
	return strings.Join(h.Groups, ", ")
}

func (h *Host) hasGroups() bool {
	return h.Groups != nil && len(h.Groups) > 0
}

func NewHostByFileName(fileName string) *Host {
	groups := strings.Split(fileName, config.SepGroupInFile)
	name := groups[len(groups)-1]
	groups = groups[:len(groups)-1]
	return &Host{
		Name:     name[:len(name)-len(config.HostFileExt)],
		FileName: fileName,
		FilePath: filepath.Join(config.BaseDir, fileName),
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

func HostsToStr(hosts []*Host) string {
	if len(hosts) == 0 {
		return ""
	}
	var hsb strings.Builder
	for _, host := range hosts {
		hsb.WriteString(host.Name)
		hsb.WriteString(", ")
	}
	return hsb.String()[:hsb.Len()-2]
}
