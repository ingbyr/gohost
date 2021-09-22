package conf

import (
	"os/user"
	"path/filepath"
)

const (
	SepGroupInFile   = "_"
	SepInCmd         = ","
	TmpCombinedHost  = ".tmp_combined"
	BaseHostFileName = "base"
)

var (
	BaseDir      string
	BaseHostFile string
)

func init() {
	curr, err := user.Current()
	if err != nil {
		panic(err)
	}

	BaseDir = filepath.Join(curr.HomeDir, ".gohost")
	BaseHostFile = filepath.Join(BaseDir, "."+BaseHostFileName)
}
