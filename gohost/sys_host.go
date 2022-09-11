package gohost

import (
	"gohost/db"
	"os"
	"sync"
)

var (
	_               Host = (*SysHost)(nil)
	sysHostOnce     sync.Once
	sysHostInstance *SysHost
)

func SysHostInstance() Host {
	sysHostOnce.Do(func() {
		sysHostInstance = &SysHost{
			Name: "System Host",
			Desc: cfg.SysHostFile,
		}
	})
	return sysHostInstance
}

type SysHost struct {
	ID   db.ID `boltholdKey:"ID"`
	Name string
	Desc string
}

func (s *SysHost) Title() string {
	return s.Name
}

func (s *SysHost) Description() string {
	return s.Desc
}

func (s *SysHost) IsEditable() bool {
	return false
}

func (s *SysHost) FilterValue() string {
	return s.Name
}

func (s *SysHost) GetParentID() db.ID {
	return 0
}

func (s *SysHost) GetID() db.ID {
	return 1
}

func (s *SysHost) GetName() string {
	return s.Name
}

func (s *SysHost) GetContent() []byte {
	hosts, err := os.ReadFile(cfg.SysHostFile)
	if err != nil {
		return []byte("Can not open system hosts file: \n" + err.Error())
	}
	return hosts
}

func (s *SysHost) SetContent(content []byte) {
	panic("implement me")
}

func (s *SysHost) GetDesc() string {
	return s.Desc
}

func (s *SysHost) IsEnabled() bool {
	return true
}
