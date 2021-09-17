/*
 @Author: ingbyr
*/

package myfs

import (
	"errors"
	"io/fs"
	"os"
)

var (
	Perm664        = fs.FileMode(0644)
	ErrNotDir      = errors.New("not a directory")
	ErrIsDir       = errors.New("is a directory")
	ErrNotEmptyDir = errors.New("not an empty directory")
)

type HostFs interface {
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
	WriteFile(path string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Stat(path string) (fs.FileInfo, error)
	IsNotExist(err error) bool
	Remove(path string) error
	Rename(oldPath, newPath string) error
}
