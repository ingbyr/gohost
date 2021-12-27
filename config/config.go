package config

import (
	"fmt"
	"github.com/ingbyr/gohost/hfs"
	"gopkg.in/ini.v1"
	"os/user"
	"path/filepath"
	"reflect"
)

const (
	Version          = "0.1.2"
	SepGroupInFile   = "_"
	SepInCmd         = ","
	BaseHostFileName = "base"
	HostFileExt      = ".txt"
)

var (
	currUser, _  = user.Current()
	BaseDir      = filepath.Join(currUser.HomeDir, ".gohost")
	BaseHostFile = filepath.Join(BaseDir, "."+BaseHostFileName)
	File         = filepath.Join(BaseDir, "config.ini")
	fs           = hfs.H
)

type Config struct {
	Editor string `ini:"editor"`
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

func (c *Config) Change(option string, value string) (error, []string) {
	var validOptions []string
	modified := false
	cfg := reflect.ValueOf(c).Elem()
	for i := 0; i < cfg.NumField(); i++ {
		field := cfg.Type().Field(i)
		op, ok := field.Tag.Lookup("ini")
		if ok {
			if op == option {
				cfg.Field(i).Set(reflect.ValueOf(value))
				modified = true
			}
			validOptions = append(validOptions, op)
		}
	}
	if !modified {
		return fmt.Errorf("not valid config (option=%s, value=%s)", option, value), validOptions
	}
	return c.Sync(), validOptions
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
