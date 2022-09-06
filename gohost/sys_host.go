package gohost

import (
	"gohost/config"
	"os"
	"sync"
)

type sysHost struct {
	name string
	desc string
}

var (
	sysHostOnce     sync.Once
	sysHostInstance *sysHost
)

func SysHost() Host {
	sysHostOnce.Do(func() {
		sysHostInstance = &sysHost{
			name: "System Host",
			desc: config.Instance().SysHostFile,
		}
	})
	return sysHostInstance
}

// Implement of Host
var _ Host = (*sysHost)(nil)

func (s *sysHost) FilterValue() string {
	return s.name
}

func (s *sysHost) GetParentID() string {
	panic("implement me")
}

func (s *sysHost) GetID() string {
	panic("implement me")
}

func (s *sysHost) GetName() string {
	//TODO implement me
	return s.name
}

func (s *sysHost) GetContent() []byte {
	hosts, err := os.ReadFile(config.Instance().SysHostFile)
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

func (s *sysHost) GetGroupID() string {
	panic("implement me")
}

func (s *sysHost) IsEnabled() bool {
	return true
}
