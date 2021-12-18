/*
 @Author: ingbyr
*/

package host

import (
	"path"
	"strings"

	"github.com/ingbyr/gohost/conf"
	"github.com/ingbyr/gohost/util"
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
	if len(groups) == 0 && hostName != conf.TmpCombinedHost {
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
