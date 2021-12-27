package config

import (
	"fmt"
	"github.com/ingbyr/gohost/hfs"
	"gopkg.in/ini.v1"
	"os/user"
	"path/filepath"
)

const (
	Version          = "0.1.2"
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

func (c *Config) Init() error {
	// create base dir
	if _, err := fs.Stat(BaseDir); fs.IsNotExist(err) {
		if err := fs.MkdirAll(BaseDir, hfs.Perm644); err != nil {
			return err
		}
	}
	// init config file or load config from file
	if _, err := fs.Stat(File); fs.IsNotExist(err) {
		newConfigFile, err := fs.Create(File)
		if err != nil {
			return err
		}
		cfg := ini.Empty()
		if err := ini.ReflectFrom(cfg, c); err != nil {
			return err
		}
		if _, err := cfg.WriteTo(newConfigFile); err != nil {
			return err
		}
		if err := newConfigFile.Close(); err != nil {
			return err
		}
	}
	file, err := hfs.H.Open(File)
	if err != nil {
		return err
	}
	cfg, err := ini.Load(file)
	if err != nil {
		return err
	}
	if err := cfg.MapTo(c); err != nil {
		return err
	}
	return nil
}

func (c *Config) Change(option string, value string) error {
	if option == "" || value == "" {
		return fmt.Errorf("not valid config (option=%s, value=%s)", option, value)
	}
	switch option {
	case OpEditor:
		c.Editor = value
	default:
		return fmt.Errorf("not valid config (option=%s, value=%s)", option, value)
	}
	return c.Sync()
}

func (c *Config) Sync() error {
	cfg := ini.Empty()
	if err := ini.ReflectFrom(cfg, c); err != nil {
		return err
	}
	cfgFile, err := hfs.H.Create(File)
	if err != nil {
		return err
	}
	if _, err = cfg.WriteTo(cfgFile); err != nil {
		return err
	}
	return nil
}
