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
	WriteFile(name string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
}
