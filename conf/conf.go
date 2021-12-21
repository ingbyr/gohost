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
	SepGroupInFile   = "_"
	SepInCmd         = ","
	TmpCombinedHost  = ".tmp_combined"
	BaseHostFileName = "base"
	HostFileExt      = ".txt"
	ModeStorage      = "storage"
	ModeMemory       = "memory"
)

type conf struct {
	Mode   string
	Editor string
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

func Sync() {
	switch C.Mode {
	case ModeStorage:
		cfg := ini.Empty()
		if err := ini.ReflectFrom(cfg, C); err != nil {
			display.ErrExit(err)
		}
		if err := cfg.SaveTo(ConfigFile); err != nil {
			display.ErrExit(err)
		}
	case ModeMemory:
	}
}
