/*
 @Author: ingbyr
*/

package hfs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
)

var (
	Perm644 = fs.FileMode(0644)

	ErrNotDir      = errors.New("not a directory")
	ErrIsDir       = errors.New("is a directory")
	ErrNotEmptyDir = errors.New("not an empty directory")
	ErrClosedFile  = errors.New("file closed")
)

type Hfs interface {
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
	WriteFile(path string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Stat(path string) (fs.FileInfo, error)
	IsNotExist(err error) bool
	Remove(path string) error
	Rename(oldPath, newPath string) error
	Create(path string) (io.WriteCloser, error)
	NewFs() Hfs
}

func printEntryTree(hfs Hfs) {
	err := fs.WalkDir(hfs, "/", func(path string, entry fs.DirEntry, err error) error {
		fmt.Printf("path %s, name %s, dir %t\n", path, entry.Name(), entry.IsDir())
		return err
	})
	if err != nil {
		panic(err)
	}
}
