//go:build linux || darwin

package config

import (
	"os"
	"path/filepath"
)

func New() *config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	baseDir := filepath.Join(homeDir, ".local", "share", name)
	if err = os.MkdirAll(baseDir, os.ModePerm); err != nil {
		panic(err)
	}
	return &config{
		BaseDir:     baseDir,
		DBFile:      filepath.Join(baseDir, name+".db"),
		SysHostFile: "/etc/hosts",
		LineBreak:   "\n",
	}
}
