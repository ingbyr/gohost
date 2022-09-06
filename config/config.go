package config

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type config struct {
	BaseDir     string
	DBFile      string
	SysHostFile string
}

func New() *config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	baseDir := filepath.Join(homeDir, ".gohost")
	_, err = os.Stat(baseDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err = os.MkdirAll(baseDir, os.ModePerm); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	dbFile := filepath.Join(baseDir, "gohost.db")
	sysHostFile := "/etc/hosts"
	return &config{
		BaseDir:     baseDir,
		DBFile:      dbFile,
		SysHostFile: sysHostFile,
	}
}

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
