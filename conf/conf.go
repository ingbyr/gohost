package conf

import (
	"fmt"
	"github.com/ingbyr/gohost/display"
	"github.com/ingbyr/gohost/editor"
	"github.com/ingbyr/gohost/hfs"
	"gopkg.in/ini.v1"
	"os/user"
	"path"
)

const (
	SepGroupInFile = "_"
	SepInCmd       = ","

	BaseHostFileName = "base"
	HostFileExt      = ".txt"
	OpEditor         = "editor"
)

type CustomConfig struct {
	Editor string `init:"editor"`
}

var (
	currUser, _  = user.Current()
	BaseDir      = path.Join(currUser.HomeDir, ".gohost")
	BaseHostFile = path.Join(BaseDir, "."+BaseHostFileName)
	ConfigFile   = path.Join(BaseDir, "config.ini")
	Custom       = &CustomConfig{
		Editor: editor.Default,
	}

	fs = hfs.H
)

func init() {
	// create base dir
	if _, err := fs.Stat(BaseDir); fs.IsNotExist(err) {
		if err := fs.MkdirAll(BaseDir, hfs.Perm644); err != nil {
			display.Panic("can not create dir "+BaseDir, err)
		}
	}
	// init config file
	if _, err := fs.Stat(ConfigFile); fs.IsNotExist(err) {
		newConfigFile, err := fs.Create(ConfigFile)
		if err != nil {
			display.ErrExit(err)
		}
		cfg := ini.Empty()
		if err := ini.ReflectFrom(cfg, Custom); err != nil {
			display.ErrExit(err)
		}
		if _, err := cfg.WriteTo(newConfigFile); err != nil {
			display.ErrExit(err)
		}
		if err := newConfigFile.Close(); err != nil {
			display.ErrExit(err)
		}
	}
	file, err := hfs.H.Open(ConfigFile)
	if err != nil {
		display.ErrExit(err)
	}
	cfg, err := ini.Load(file)
	if err != nil {
		display.ErrExit(err)
	}
	if err := cfg.MapTo(Custom); err != nil {
		display.ErrExit(err)
	}
}

func Change(op string, value string) {
	if op == "" || value == "" {
		display.Err(fmt.Errorf("not valid config (op=%s, value=%s)", op, value))
	}
	switch op {
	case OpEditor:
		Custom.Editor = value
	default:
		display.ErrExit(fmt.Errorf("not valid config (op=%s, value=%s)", op, value))
	}
	Sync()
}

func Sync() {
	cfg := ini.Empty()
	if err := ini.ReflectFrom(cfg, Custom); err != nil {
		display.ErrExit(err)
	}

	// todo writer interface
	if err := cfg.SaveTo(ConfigFile); err != nil {
		display.ErrExit(err)
	}
}
