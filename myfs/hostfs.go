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
	ErrClosedFile  = errors.New("file closed")
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

//func PrintEntryTree(hostFs HostFs)  {
//	err := fs.WalkDir(hostFs, rootDirPath, func(path string, entry fs.DirEntry, err error) error {
//		fmt.Println(path)
//		fmt.Println(entry.Name())
//		return err
//	})
//	if err != nil {
//		panic(err)
//	}
//}