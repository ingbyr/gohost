package main

import (
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	BaseDir string
	DBFile  string
}

var cfg *Config

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	baseDir := filepath.Join(homeDir, ".local", "share", "gohost")
	_, err = os.Stat(baseDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(filepath.Dir(baseDir), os.ModePerm); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	dbFile := filepath.Join(baseDir, "gohost.db")
	_, err = os.Stat(dbFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	cfg = &Config{
		BaseDir: baseDir,
		DBFile:  dbFile,
	}
}
