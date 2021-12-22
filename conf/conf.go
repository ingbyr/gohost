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

	OpMode   = "mode"
	OpEditor = "editor"
)

type conf struct {
	Mode   string `ini:"mode"`
	Editor string `init:"editor"`
}

var (
	currUser, _  = user.Current()
	BaseDir      = path.Join(currUser.HomeDir, ".gohost")
	BaseHostFile = path.Join(BaseDir, "."+BaseHostFileName)
	ConfigFile   = path.Join(BaseDir, "config.ini")
	C            = &conf{
		Mode:   ModeStorage,
		Editor: editor.Default,
	}
)

func init() {
	// init config file
	// todo reader interface
	_ = ini.MapTo(C, ConfigFile)
	checkMode()
}

func checkMode() {
	switch C.Mode {
	case ModeStorage, ModeMemory:
	default:
		display.ErrExit(fmt.Errorf("not valid mode %s", C.Mode))
	}
}

func Change(op string, value string) {
	if op == "" || value == "" {
		display.Err(fmt.Errorf("not valid config (op=%s, value=%s)", op, value))
	}
	switch op {
	case OpMode:
		C.Mode = value
		checkMode()
	case OpEditor:
		C.Editor = value
	default:
		display.ErrExit(fmt.Errorf("not valid config (op=%s, value=%s)", op, value))
	}
	Sync()
}

func Sync() {
	switch C.Mode {
	case ModeStorage:
		cfg := ini.Empty()
		if err := ini.ReflectFrom(cfg, C); err != nil {
			display.ErrExit(err)
		}

		// todo writer interface
		if err := cfg.SaveTo(ConfigFile); err != nil {
			display.ErrExit(err)
		}
	case ModeMemory:
	}
}
