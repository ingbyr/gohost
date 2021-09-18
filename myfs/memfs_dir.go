/*
 @Author: ingbyr
*/

package myfs

import (
	"errors"
	"io/fs"
	"sort"
	"time"
)

type MemDir struct {
	name     string
	modTime  time.Time
	children map[string]fs.DirEntry
	mode     fs.FileMode
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

	entries := make([]fs.DirEntry, 0, n)
	for i := 0; i < n && i < totalEntry; i++ {
		name := names[i]
		child := m.children[name]

		f, isFile := child.(*MemFile)
		if isFile {
			entries = append(entries, f)
		} else {
			dirEntry := child.(*MemDir)
			entries = append(entries, dirEntry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	return entries, nil
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
