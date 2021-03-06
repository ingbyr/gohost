/*
 @Author: ingbyr
*/

package hfs

import (
	"bytes"
	"io/fs"
	"os"
	"time"
)

type MemFile struct {
	name    string
	content *bytes.Buffer
	modTime time.Time
	mode    os.FileMode
	closed  bool
}

func (m *MemFile) Write(p []byte) (n int, err error) {
	return m.content.Write(p)
}

func (m *MemFile) Stat() (fs.FileInfo, error) {
	return m, nil
}

func (m *MemFile) Read(bytes []byte) (int, error) {
	if m.closed {
		return 0, ErrClosedFile
	}
	return m.content.Read(bytes)
}

func (m *MemFile) Close() error {
	m.closed = true
	return nil
}

func (m *MemFile) Name() string {
	return m.name
}

func (m *MemFile) Size() int64 {
	return int64(m.content.Len())
}

func (m *MemFile) Mode() fs.FileMode {
	return m.mode
}

func (m *MemFile) ModTime() time.Time {
	return m.modTime
}

func (m *MemFile) IsDir() bool {
	return false
}

func (m *MemFile) Sys() interface{} {
	return nil
}

func (m *MemFile) Type() fs.FileMode {
	return m.Mode()
}

func (m *MemFile) Info() (fs.FileInfo, error) {
	return m.Stat()
}
