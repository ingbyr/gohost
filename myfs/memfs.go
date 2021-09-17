/*
 @Author: ingbyr
*/

package myfs

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	_ HostFs = NewMemFs()
)

type MemFs struct {
	rootDir *MemDir
}

func NewMemFs() *MemFs {
	return &MemFs{
		rootDir: &MemDir{
			children: make(map[string]fs.DirEntry),
		},
	}
}

func (m *MemFs) Open(path string) (fs.File, error) {
	path = validPath(path)
	if path == invalidPath {
		return nil, &fs.PathError{
			Op:   "open",
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}

	cur := m.rootDir
	parts := strings.Split(path, "/")
	for i, part := range parts {
		child := cur.children[part]
		if child == nil {
			return nil, &fs.PathError{
				Op:   "open",
				Path: path,
				Err:  fs.ErrNotExist,
			}
		}

		f, ok := child.(*MemFile)
		if ok {
			if i == len(parts)-1 {
				f.closed = false
				return f, nil
			}
			return nil, &fs.PathError{
				Op:   "open",
				Path: path,
				Err:  fs.ErrNotExist,
			}
		}

		d, ok := child.(*MemDir)
		if !ok {
			return nil, &fs.PathError{
				Op:   "open",
				Path: path,
				Err:  ErrNotDir,
			}
		}

		cur = d
	}

	return cur, nil
}

func (m *MemFs) ReadDir(path string) ([]fs.DirEntry, error) {
	path = validPath(path)
	if path == invalidPath {
		return nil, &fs.PathError{
			Op:   "write",
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}
	dir, err := m.getDir(path)
	if err != nil {
		return nil, err
	}
	return dir.ReadDir(0)
}

func (m *MemFs) ReadFile(path string) ([]byte, error) {
	f, err := m.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	stat, err := f.Stat()
	var result = make([]byte, stat.Size())
	if _, err := f.Read(result); err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MemFs) WriteFile(path string, data []byte, perm os.FileMode) error {
	path = validPath(path)
	if path == invalidPath {
		return &fs.PathError{
			Op:   "write",
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}

	var err error
	dir := m.rootDir
	parentPath := filepath.Dir(path)
	if parentPath != "." {
		dir, err = m.getDir(parentPath)
		if err != nil {
			return err
		}
	}
	filename := filepath.Base(path)
	dir.children[filename] = &MemFile{
		name:    filename,
		content: bytes.NewBuffer(data),
		modTime: time.Now(),
		mode:    perm,
		closed:  true,
	}
	return nil
}

func (m *MemFs) MkdirAll(path string, perm os.FileMode) error {
	path = validPath(path)
	if path == invalidPath {
		return &fs.PathError{
			Op:   "MkdirAll",
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}
	cur := m.rootDir
	parts := strings.Split(path, "/")
	for _, part := range parts {
		child, ok := cur.children[part]
		if !ok {
			childDir := &MemDir{
				name:     part,
				modTime:  time.Now(),
				children: make(map[string]fs.DirEntry),
				mode:     fs.ModeDir | perm,
			}
			cur.children[part] = childDir
			cur = childDir
		} else {
			childDir, ok := child.(*MemDir)
			if !ok {
				return &fs.PathError{
					Op:   "MkdirAll",
					Path: path,
					Err:  ErrNotDir,
				}
			}
			cur = childDir
		}
	}
	return nil
}

func (m *MemFs) Stat(path string) (fs.FileInfo, error) {
	path = validPath(path)
	if path == invalidPath {
		return nil, &fs.PathError{
			Op:   "Stat",
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}

	dirPath := filepath.Dir(path)
	dir, err := m.getDir(dirPath)
	if err != nil {
		return nil, err
	}
	dirEntry, ok := dir.children[filepath.Base(path)]
	if !ok {
		return nil, &fs.PathError{
			Op:   "Stat",
			Path: path,
			Err:  fs.ErrNotExist,
		}
	}
	if dirEntry.IsDir() {
		return dirEntry.(*MemDir).Stat()
	}
	return dirEntry.(*MemFile).Stat()
}

func (m *MemFs) IsNotExist(err error) bool {
	fsPathError, ok := err.(*fs.PathError)
	if !ok {
		return false
	}
	return fsPathError.Err == fs.ErrNotExist
}

func (m *MemFs) Remove(path string) error {
	_path := validPath(path)
	if _path == invalidPath {
		return &fs.PathError{
			Op:   "Remove",
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}

	parentDir, err := m.getDir(filepath.Dir(_path))
	if err != nil {
		return err
	}
	name := filepath.Base(_path)
	if entry, ok := parentDir.children[name]; ok {
		// if target entry has children return ErrNotEmptyDir error
		if entry.IsDir() && len(entry.(*MemDir).children) > 0 {
			return &fs.PathError{
				Op:   "Remove",
				Path: path,
				Err:  ErrNotEmptyDir,
			}
		}
		delete(parentDir.children, name)
	}
	return nil
}

func (m *MemFs) Rename(oldPath, newPath string) error {
	panic("implement me")
}

// getFile path must be valid by caller
func (m *MemFs) getFile(path string) (*MemFile, error) {
	entry, err := m.getEntry(path)
	if err != nil {
		return nil, err
	}
	if entry.IsDir() {
		return nil, &fs.PathError{
			Op:   "getFile",
			Path: path,
			Err:  ErrIsDir,
		}
	}
	return entry.(*MemFile), nil
}

// getDir path must be valid by caller
func (m *MemFs) getDir(path string) (*MemDir, error) {
	entry, err := m.getEntry(path)
	if err != nil {
		return nil, err
	}
	if !entry.IsDir() {
		return nil, &fs.PathError{
			Op:   "getDir",
			Path: path,
			Err:  ErrNotDir,
		}
	}
	return entry.(*MemDir), nil
}

// getEntry path must be valid by caller
func (m *MemFs) getEntry(path string) (fs.DirEntry, error) {
	parts := strings.Split(path, "/")
	dir := m.rootDir
	var targetEntry fs.DirEntry
	parentDirLen := len(parts) - 1
	for i, part := range parts {
		entry := dir.children[part]
		if entry == nil {
			return nil, &fs.PathError{
				Op:   "getEntry",
				Path: path,
				Err:  fs.ErrNotExist,
			}
		}
		if i < parentDirLen {
			// getDir is in parent getDir paths
			if !entry.IsDir() {
				return nil, &fs.PathError{
					Op:   "getDir",
					Path: path,
					Err:  ErrNotDir,
				}
			}
			dir = entry.(*MemDir)
		} else {
			// find target
			targetEntry = entry
		}
	}
	return targetEntry, nil
}
