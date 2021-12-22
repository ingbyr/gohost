package conf

import (
	"fmt"
	"github.com/ingbyr/gohost/display"
	"github.com/ingbyr/gohost/editor"
	"gopkg.in/ini.v1"
	"os/user"
	"path"
)

const (
	SepGroupInFile = "_"
	SepInCmd       = ","

	TmpCombinedHost  = ".tmp_combined"
	BaseHostFileName = "base"
	HostFileExt      = ".txt"

	ModeStorage = "storage"
	ModeMemory  = "memory"

	OpEditor = "editor"
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
)

func init() {
	// init config file
	// todo reader interface

	_ = ini.MapTo(Custom, ConfigFile)
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
