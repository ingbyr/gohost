/*
 @Author: ingbyr
*/

package fss

import (
	"io/fs"
	"os"
)

var (
	_ HostFs = NewOsFs()
)

type OsFs struct {
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
