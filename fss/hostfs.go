/*
 @Author: ingbyr
*/

package fss

import (
	"errors"
	"io/fs"
	"os"
)

var (
	ErrNotDir = errors.New("not a directory")
)

type HostFs interface {
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
	WriteFile(name string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Stat(name string) (fs.FileInfo, error)
	IsNotExist(err error) bool
}
