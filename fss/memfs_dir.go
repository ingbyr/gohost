/*
 @Author: ingbyr
*/

package fss

import (
	"errors"
	"io/fs"
	"time"
)

type MemDir struct {
	name string
	modTime time.Time
	children map[string]fs.DirEntry
	// idx file index in current dir
	idx int
	mode fs.FileMode
}

func (m *MemDir) Read(p []byte) (int, error) {
	return 0, &fs.PathError{
		Op:   "read",
		Path: m.name,
		Err:  errors.New("is directory"),
	}
}

func (m *MemDir) Stat() (fs.FileInfo, error) {
	return m, nil
}

func (m *MemDir) Close() error {
	return nil
}

func (m *MemDir) ReadDir(n int) ([]fs.DirEntry, error) {
	names := make([]string, 0, len(m.children))
	for name := range m.children {
		names = append(names, name)
	}

	totalEntry := len(names)
	if n <= 0 {
		n = totalEntry
	}

	dirEntries := make([]fs.DirEntry, 0, n)
	for i := m.idx; i < n && i < totalEntry; i++ {
		name := names[i]
		child := m.children[name]

		f, isFile := child.(*MemFile)
		if isFile {
			dirEntries = append(dirEntries, f)
		} else {
			dirEntry := child.(*MemDir)
			dirEntries = append(dirEntries, dirEntry)
		}

		m.idx = i
	}

	return dirEntries, nil
}

func (m *MemDir) Name() string {
	return m.name
}

func (m *MemDir) Size() int64 {
	return 0
}

func (m *MemDir) Mode() fs.FileMode {
	return m.mode
}

func (m *MemDir) ModTime() time.Time {
	return m.modTime
}

func (m *MemDir) IsDir() bool {
	return true
}

func (m *MemDir) Sys() interface{} {
	return nil
}

func (m *MemDir) Type() fs.FileMode {
	return m.Mode()
}

func (m *MemDir) Info() (fs.FileInfo, error) {
	return m.Stat()
}
