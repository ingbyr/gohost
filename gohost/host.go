package gohost

import (
	"bytes"
	"os"
)

type Host interface {
	Node
	GetContent() []byte
	SetContent([]byte)
	IsEditable() bool
}

func (s *Service) ApplyHost(hostContent []byte) {
	// Truncate system host file
	sysHostFile, err := os.Create(cfg.SysHostFile)
	if err != nil {
		panic(err)
	}
	defer sysHostFile.Close()
	// Write hosts to system host file
	if _, err = sysHostFile.Write(hostContent); err != nil {
		panic(err)
	}
}

func (s *Service) CombineEnabledHosts() []byte {
	hosts := s.loadLocalHostsByFlag(MaskEnable)
	// TODO load all enabled remote hosts
	combinedHost := bytes.NewBuffer(nil)
	for _, host := range hosts {
		combinedHost.WriteString("# Content from ")
		combinedHost.WriteString(host.Title())
		if host.Description() != "" {
			combinedHost.WriteString("( ")
			combinedHost.WriteString(host.Description())
			combinedHost.WriteString(" )")
		}
		combinedHost.WriteString(cfg.LineBreak)
		combinedHost.Write(host.GetContent())
		combinedHost.WriteString(cfg.LineBreak)
		combinedHost.WriteString("# End of ")
		combinedHost.WriteString(host.Title())
		combinedHost.WriteString(cfg.LineBreak)
		combinedHost.WriteString(cfg.LineBreak)
	}
	return combinedHost.Bytes()
}
