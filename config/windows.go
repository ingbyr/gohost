//go:build windows

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
	baseDir := filepath.Join(homeDir, name)
	if err = os.MkdirAll(baseDir, os.ModePerm); err != nil {
		panic(err)
	}
	return &config{
		BaseDir:     baseDir,
		DBFile:      filepath.Join(baseDir, "gohost.db"),
		SysHostFile: "C:\\Windows\\System32\\drivers\\etc\\hosts",
	}
}
