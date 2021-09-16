/*
 @Author: ingbyr
*/

package fss

import (
	"io/fs"
	"os"
)

type DiskFs struct {
}

func (d *DiskFs) Open(path string) (fs.File, error) {
	return os.Open(path)
}

func (d *DiskFs) ReadDir(path string) ([]fs.DirEntry, error) {
	return os.ReadDir(path)
}

func (d *DiskFs) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (d *DiskFs) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (d *DiskFs) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
