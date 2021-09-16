/*
 @Author: ingbyr
*/

package fss

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	if path == "" {
		return nil, &fs.PathError{
			Op:   "open",
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}

	if path == "/" || path == "" {
		m.rootDir.idx = 0
		return m.rootDir, nil
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
				Err:  errors.New("not a directory"),
			}
		}

		d.idx = 0
		cur = d
	}

	return cur, nil
}

func (m *MemFs) ReadDir(path string) ([]fs.DirEntry, error) {
	path = validPath(path)
	if path == "" {
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
	if path == "" {
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
	if path == "" {
		return &fs.PathError{
			Op:   "mkdir",
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
				mode:     perm,
			}
			cur.children[part] = childDir
			cur = childDir
		} else {
			childDir, ok := child.(*MemDir)
			if !ok {
				return fmt.Errorf("%s is not directory", part)
			}
			cur = childDir
		}
	}
	return nil
}

// getDir path is already validated by caller
func (m *MemFs) getDir(path string) (*MemDir, error) {
	parts := strings.Split(path, "/")
	cur := m.rootDir
	for _, part := range parts {
		child := cur.children[part]
		if child == nil {
			return nil, fmt.Errorf("%s is not exists", path)
		}
		childDir, ok := child.(*MemDir)
		if !ok {
			return nil, fmt.Errorf("%s is not directory", path)
		}
		cur = childDir
	}

	return cur, nil
}
