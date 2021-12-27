//go:build osfs
// +build osfs

/*
 @Author: ingbyr
*/

package hfs

import (
	"io"
	"io/fs"
	"os"
)

var H Hfs = NewOsFs()

type OsFs struct {
}

func (o *OsFs) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (o *OsFs) NewFs() Hfs {
	panic("os fs do not support NewFs() method")
}

func (o *OsFs) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func NewOsFs() *OsFs {
	return &OsFs{}
}

func (o *OsFs) Open(path string) (fs.File, error) {
	return os.Open(path)
}

func (o *OsFs) ReadDir(path string) ([]fs.DirEntry, error) {
	return os.ReadDir(path)
}

func (o *OsFs) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (o *OsFs) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (o *OsFs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (o *OsFs) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (o *OsFs) Remove(name string) error {
	return os.Remove(name)
}

func (o *OsFs) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func (o *OsFs) Create(path string) (io.WriteCloser, error) {
	return os.Create(path)
}
