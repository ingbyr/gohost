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
			cfg.SysHostFile = "hosts"
		}
	})
	return cfg
}
