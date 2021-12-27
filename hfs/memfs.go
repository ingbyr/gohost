//go:build memfs
// +build memfs

/*
 @Author: ingbyr
*/

package hfs

import (
	"bytes"
	"github.com/ingbyr/gohost/display"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var H Hfs = NewMemFs()

const (
	rootDirName = "mem"
)

type MemFs struct {
	rootDir *MemDir
}

func (m *MemFs) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	//TODO implement me
	panic("mem fs not support OpenFile method")
}

func (m *MemFs) NewFs() Hfs {
	return &MemFs{
		rootDir: &MemDir{
			children: make(map[string]fs.DirEntry),
		},
	}
}

func (m *MemFs) Reset() {
}

func NewMemFs() *MemFs {
	display.Warn("memory mode")
	return &MemFs{
		rootDir: &MemDir{
			children: make(map[string]fs.DirEntry),
		},
	}
}

func (m *MemFs) Open(path string) (fs.File, error) {
	path, ok := m.validPath(path)
	if !ok {
		return nil, &fs.PathError{
			Op:   "open",
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}

	cur := m.rootDir
	parts := strings.Split(path, string(filepath.Separator))
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
	path, ok := m.validPath(path)
	if !ok {
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
	path, ok := m.validPath(path)
	if !ok {
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
	path, ok := m.validPath(path)
	if !ok {
		return &fs.PathError{
			Op:   "MkdirAll",
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}
	cur := m.rootDir
	parts := strings.Split(path, string(filepath.Separator))
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
	path, ok := m.validPath(path)
	if !ok {
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
	path, ok := m.validPath(path)
	if !ok {
		return &fs.PathError{
			Op:   "Remove",
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}

	parentDir, err := m.getDir(filepath.Dir(path))
	if err != nil {
		return err
	}
	name := filepath.Base(path)
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

// Rename only support rename file name
func (m *MemFs) Rename(oldPath, newPath string) error {
	oldPath, ok := m.validPath(oldPath)
	if !ok {
		return &fs.PathError{
			Op:   "Rename",
			Path: oldPath,
			Err:  fs.ErrInvalid,
		}
	}
	newPath, ok = m.validPath(newPath)
	if !ok {
		return &fs.PathError{
			Op:   "Rename",
			Path: oldPath,
			Err:  fs.ErrInvalid,
		}
	}

	oldParentDir, err := m.getDir(filepath.Dir(oldPath))
	if err != nil {
		return err
	}
	oldFileName := filepath.Base(oldPath)
	oldEntry, ok := oldParentDir.children[oldFileName]
	if !ok {
		return &fs.PathError{
			Op:   "Rename",
			Path: oldPath,
			Err:  fs.ErrNotExist,
		}
	}
	if oldEntry.IsDir() {
		return &fs.PathError{
			Op:   "Rename",
			Path: oldPath,
			Err:  ErrIsDir,
		}
	}
	oldFile := oldEntry.(*MemFile)

	newParent, err := m.getDir(filepath.Dir(newPath))
	if err != nil {
		return err
	}
	newFileName := filepath.Base(newPath)
	newEntry, ok := newParent.children[newFileName]
	if ok && newEntry.IsDir() {
		return &fs.PathError{
			Op:   "Rename",
			Path: newPath,
			Err:  ErrIsDir,
		}
	}
	oldFile.name = newFileName
	newParent.children[newFileName] = oldFile
	delete(oldParentDir.children, oldFileName)
	return nil
}

func (m *MemFs) Create(path string) (io.WriteCloser, error) {
	path, ok := m.validPath(path)
	if !ok {
		return nil, &fs.PathError{
			Op:   "Create",
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}

	dirPath := filepath.Dir(path)
	dir, err := m.getDir(dirPath)
	if err != nil {
		return nil, err
	}
	fileName := filepath.Base(path)
	file := &MemFile{
		name:    fileName,
		content: bytes.NewBuffer([]byte("")),
		modTime: time.Now(),
		mode:    Perm644,
		closed:  false,
	}
	dir.children[fileName] = file
	return file, nil
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
	if path == "." {
		return m.rootDir, nil
	}
	parts := strings.Split(path, string(filepath.Separator))
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

func (m *MemFs) validPath(path string) (string, bool) {
	memPath := filepath.Join(rootDirName, path)
	if fs.ValidPath(memPath) {
		return memPath, true
	}
	return "", false
}
