package config

import (
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
	cfg  *config
	once sync.Once
)

func Instance() *config {
	once.Do(func() {
		cfg = New()
	})
	return cfg
}
