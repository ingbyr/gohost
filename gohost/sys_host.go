package gohost

import (
	"os"
	"sync"
)

var (
	_               Host = (*sysHost)(nil)
	sysHostOnce     sync.Once
	sysHostInstance *sysHost
)

func SysHost() Host {
	sysHostOnce.Do(func() {
		sysHostInstance = &sysHost{
			name: "System Host",
			desc: cfg.SysHostFile,
		}
	})
	return sysHostInstance
}

type sysHost struct {
	name string
	desc string
}

func (s *sysHost) Title() string {
	return s.name
}

func (s *sysHost) Description() string {
	return s.desc
}

func (s *sysHost) IsEditable() bool {
	return false
}

func (s *sysHost) FilterValue() string {
	return s.name
}

func (s *sysHost) GetParentID() string {
	return ""
}

func (s *sysHost) GetID() string {
	return "-1"
}

func (s *sysHost) GetName() string {
	return s.name
}

func (s *sysHost) GetContent() []byte {
	hosts, err := os.ReadFile(cfg.SysHostFile)
	if err != nil {
		return []byte("Can not open system hosts file: \n" + err.Error())
	}
	return hosts
}

func (s *sysHost) SetContent(content []byte) {
	panic("implement me")
}

func (s *sysHost) GetDesc() string {
	return s.desc
}

func (s *sysHost) IsEnabled() bool {
	return true
}
