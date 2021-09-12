package _const

import (
	"os/user"
	"path/filepath"
)

const (
	SepGroup        = "_"
	TmpCombinedHost = ".tmp_combined"
)

var (
	BaseHostFile string
	BaseDir      string
	ConfigFile   string
)

func init() {
	curr, err := user.Current()
	if err != nil {
		panic(err)
	}
	BaseDir = filepath.Join(curr.HomeDir, ".gohost")
	BaseHostFile = filepath.Join(BaseDir, "base")
	ConfigFile = filepath.Join(BaseDir, ".conf")
}
