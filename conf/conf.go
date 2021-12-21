package conf

import (
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
)

type conf struct {
	Editor string `ini:"editor"`
}

var (
	currUser, _  = user.Current()
	BaseDir      = path.Join(currUser.HomeDir, ".gohost")
	BaseHostFile = path.Join(BaseDir, "."+BaseHostFileName)
	ConfigFile   = path.Join(BaseDir, "config.ini")
	Conf         = &conf{
		Editor: editor.Default,
	}
)

func init() {
	// init config file
	err := ini.MapTo(Conf, ConfigFile)
	if err != nil {
		Sync()
	}
}

func Sync() {
	cfg := ini.Empty()
	if err := ini.ReflectFrom(cfg, Conf); err != nil {
		display.ErrExit(err)
	}
	if err := cfg.SaveTo(ConfigFile); err != nil {
		display.ErrExit(err)
	}
}
