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
	Name     string
	FileName string
	FilePath string
	Groups   []string
}

func (h *Host) GroupsAsStr() string {
	return strings.Join(h.Groups, ", ")
}

func NewHostByFileName(fileName string) *Host {
	groups := strings.Split(fileName, conf.SepGroupInFile)
	name := groups[len(groups)-1]
	groups = groups[:len(groups)-1]
	groups = util.SortUniqueStringSlice(groups)
	return &Host{
		Name:     name,
		FileName: fileName,
		FilePath: path.Join(conf.BaseDir, fileName),
		Groups:   groups,
	}
}

func NewHostByNameGroups(hostName string, groups []string) *Host {
	if len(groups) == 0 {
		groups = append(groups, hostName)
	} else {
		groups = util.SortUniqueStringSlice(groups)
	}
	fileName := strings.Join(append(groups, hostName), conf.SepGroupInFile)
	return &Host{
		Name:     hostName,
		FileName: fileName,
		FilePath: path.Join(conf.BaseDir, fileName),
		Groups:   groups,
	}
}
