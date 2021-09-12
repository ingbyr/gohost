/*
 @Author: ingbyr
*/

package host

import (
	"github.com/ingbyr/gohost/const"
	"strings"
)

type Host struct {
	Name     string
	FileName string
	Path     string
	Groups   []string
}

func (h *Host) GroupsAsStr() string {
	return strings.Join(h.Groups, ", ")
}

func NewHost(fileName string, path string) *Host {
	groups := strings.Split(fileName, _const.SepGroup)
	name := groups[len(groups)-1]
	if len(groups) > 1 {
		groups = groups[:len(groups)-1]
	}
	return &Host{
		Name:     name,
		FileName: fileName,
		Path:     path,
		Groups:   groups,
	}
}
