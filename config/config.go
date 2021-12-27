package config

import (
	"fmt"
	"github.com/ingbyr/gohost/display"
	"github.com/ingbyr/gohost/hfs"
	"gopkg.in/ini.v1"
	"os/user"
	"path/filepath"
)

const (
	Version          = "0.1.1"
	SepGroupInFile   = "_"
	SepInCmd         = ","
	BaseHostFileName = "base"
	HostFileExt      = ".txt"
	OpEditor         = "editor"
)

var (
	currUser, _  = user.Current()
	BaseDir      = filepath.Join(currUser.HomeDir, ".gohost")
	BaseHostFile = filepath.Join(BaseDir, "."+BaseHostFileName)
	File         = filepath.Join(BaseDir, "config.ini")
	fs           = hfs.H
)

type Config struct {
	Editor     string `ini:"editor"`
	EditorArgs string `ini:"editor_args"`
}

func (c *Config) Init() {
	// create base dir
	if _, err := fs.Stat(BaseDir); fs.IsNotExist(err) {
		if err := fs.MkdirAll(BaseDir, hfs.Perm644); err != nil {
			display.Panic("can not create dir "+BaseDir, err)
		}
	}
	// init config file or load config from file
	if _, err := fs.Stat(File); fs.IsNotExist(err) {
		newConfigFile, err := fs.Create(File)
		if err != nil {
			display.ErrExit(err)
		}
		cfg := ini.Empty()
		if err := ini.ReflectFrom(cfg, c); err != nil {
			display.ErrExit(err)
		}
		if _, err := cfg.WriteTo(newConfigFile); err != nil {
			display.ErrExit(err)
		}
		if err := newConfigFile.Close(); err != nil {
			display.ErrExit(err)
		}
	}
	file, err := hfs.H.Open(File)
	if err != nil {
		display.ErrExit(err)
	}
	cfg, err := ini.Load(file)
	if err != nil {
		display.ErrExit(err)
	}
	if err := cfg.MapTo(c); err != nil {
		display.ErrExit(err)
	}
}

func (c *Config) Change(op string, value string) {
	if op == "" || value == "" {
		display.Err(fmt.Errorf("not valid config (op=%s, value=%s)", op, value))
	}
	switch op {
	case OpEditor:
		c.Editor = value
	default:
		display.ErrExit(fmt.Errorf("not valid config (op=%s, value=%s)", op, value))
	}
	c.Sync()
}

func (c *Config) Sync() {
	cfg := ini.Empty()
	if err := ini.ReflectFrom(cfg, c); err != nil {
		display.ErrExit(err)
	}
	cfgFile, err := hfs.H.Create(File)
	if err != nil {
		display.ErrExit(err)
	}
	if _, err = cfg.WriteTo(cfgFile); err != nil {
		display.ErrExit(err)
	}
}
