/*
 @Author: ingbyr
*/

package fss

import (
	"io/fs"
	"os"
)

type HostFs interface {
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
	MkdirAll(path string, perm os.FileMode) error
	WriteFile(name string, data []byte, perm os.FileMode) error
}
