package config

import (
	"os"
	"sync"
)

type config struct {
	BaseDir     string
	DBFile      string
	SysHostFile string
	LineBreak   string
}

const (
	name = "gohost"
)

var (
	cfg   *config
	once  sync.Once
	debug = os.Getenv("GOHOST_DEBUG") == "true"
)

func Instance() *config {
	once.Do(func() {
		cfg = New()
		if debug {
			cfg.SysHostFile = "fake_sys_hosts"
			_, err := os.Stat(cfg.SysHostFile)
			if err != nil {
				if os.IsNotExist(err) {
					_, err := os.Create(cfg.SysHostFile)
					if err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}
			}
		}
	})
	return cfg
}
