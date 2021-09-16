/*
 @Author: ingbyr
*/

package fss

import (
	"io/fs"
)

type DiskFs struct {

}

func (d *DiskFs) Open(name string) (fs.File, error) {
	panic("implement me")
}

func (d *DiskFs) ReadDir(name string) ([]fs.DirEntry, error) {
	panic("implement me")
}

func (d *DiskFs) ReadFile(name string) ([]byte, error) {
	panic("implement me")
}


