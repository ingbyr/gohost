package conf

import (
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

var (
	BaseDir      string
	BaseHostFile string
	ConfigFile   string
)

func init() {
	curr, err := user.Current()
	if err != nil {
		panic(err)
	}

	BaseDir = path.Join(curr.HomeDir, ".gohost")
	BaseHostFile = path.Join(BaseDir, "."+BaseHostFileName)
	ConfigFile = path.Join(BaseDir, "config.ini")
}
