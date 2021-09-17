/*
 @Author: ingbyr
*/

package myfs

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ingbyr/gohost/util"
	"io/fs"
	"path/filepath"
	"testing"
)

func TestDemo(t *testing.T) {
	diff := cmp.Diff([]string{"a", "b"}, []string{"b", "a"}, cmpopts.SortSlices(util.CmpStringSliceLess))
	fmt.Println(diff)
}

func TestMemFsUtil_ValidatePath(t *testing.T) {
	var tests = []struct {
		path string
		want string
	}{
		{"", ""},
		{"./", ""},
		{"..", ""},
		{"/", "mem"},
		{"/a/bb/ccc/d", "mem/a/bb/ccc/d"},
		{"/a/bbb/ccc/", ""},
		{"/a/bb/ccc/./d", ""},
		{"/a/bbb/../ccc", ""},
	}

	for _, test := range tests {
		if diff := cmp.Diff(test.want, validPath(test.path)); diff != "" {
			t.Fatalf("\n path %s\n diff: %s", test.path, diff)
		}
	}
}

func TestMemFs_CreateDir(t *testing.T) {
	var memFs HostFs = NewMemFs()
	dirs := []string{
		"/d1/d11",
		"/d1/d12/d111",
		"/d2",
	}
	for _, dir := range dirs {
		if err := memFs.MkdirAll(dir, fs.ModeDir|0644); err != nil {
			t.Fatal(err)
		}
	}
	files := []string{
		"/f1",
		"/f2",
		"/d1/f11",
		"/d2/f22",
	}
	for _, f := range files {
		if err := memFs.WriteFile(f, []byte("test"), 0664); err != nil {
			t.Fatal(err)
		}
	}

	var tests = []struct {
		dir    string
		wanted []string
	}{
		{"/", []string{"d1", "d2", "f1", "f2"}},
		{"/d1", []string{"d11", "d12", "f11"}},
	}

	for _, test := range tests {
		if dirs, err := memFs.ReadDir(test.dir); err != nil {
			t.Fatal(err)
		} else {
			subDirNames := make([]string, 0)
			for _, dir := range dirs {
				subDirNames = append(subDirNames, dir.Name())
			}
			if diff := cmp.Diff(subDirNames, test.wanted, cmpopts.SortSlices(util.CmpStringSliceLess)); diff != "" {
				t.Fatalf("input: %s \ndiff: %s", test.dir, diff)
			}
		}
	}

	dir := "/x/y"
	if err := memFs.MkdirAll(dir, Perm664); err != nil {
		t.Fatal(err)
	}

	if _, err := memFs.Open("/x/not_y"); !memFs.IsNotExist(err) {
		t.Fatal("should be 'not exist' err", err)
	}
	if err := memFs.WriteFile("/x/y/f", []byte("content"), Perm664); err != nil {
		t.Fatal(err)
	}
	if err := memFs.MkdirAll("/x/y/f/rush-b", Perm664); err.(*fs.PathError).Err != ErrNotDir {
		t.Fatal("should be 'not a directory' err", err)
	}
}

func TestMemFs_WriteRead(t *testing.T) {
	var tests =  []struct{
		dir string
		file string
		content []byte
	} {
		{"/", "f1", []byte("f1")},
		{"/a/b", "f2", []byte("f2")},
		{"/a/c", "f3", []byte("f3")},
	}

	var memFs HostFs = NewMemFs()
	for _, test := range tests {
		// create dirs
		if err := memFs.MkdirAll(test.dir, Perm664); err != nil {
			t.Fatal(err)
		}
		// write getFile
		filePath := filepath.Join(test.dir, test.file)
		if err := memFs.WriteFile(filePath, test.content, Perm664); err != nil {
			t.Fatal(err)
		}
		// open getFile
		file, err := memFs.Open(filePath)
		if err != nil {
			t.Fatal(err)
		}
		stat, err := file.Stat()
		if err != nil {
			t.Fatal(err)
		}
		var result = make([]byte, int(stat.Size()))
		n, err := file.Read(result)
		if err != nil {
			t.Fatal(err)
		}
		// compare
		if diff := cmp.Diff(result[:n], test.content); diff != "" {
			t.Fatalf("diff %s", diff)
		}
		if err := file.Close(); err != nil{
			t.Fatal(err)
		}
	}
}

func TestMemFs_Remove(t *testing.T) {
	memFs := NewMemFs()
	dir := "/a/b"
	if err := memFs.MkdirAll(dir, Perm664); err != nil {
		t.Fatal(err)
	}
	if err := memFs.WriteFile(filepath.Join(dir, "c1"), []byte("c1"), Perm664); err != nil {
		t.Fatal(err)
	}
	if err := memFs.WriteFile(filepath.Join(dir, "c2"), []byte("c2"), Perm664); err != nil {
		t.Fatal(err)
	}

	// remove /a/b/c1
	if err := memFs.Remove(filepath.Join(dir, "c1")); err != nil {
		t.Fatal(err)
	}

	// left /a/b/c2
	entries, _ := memFs.ReadDir(dir)
	if len(entries) != 1 || entries[0].Name() != "c2" {
		t.Fatal("error entries", entries)
	}

	// try to remove /a which has children entry
	if err := memFs.Remove("/a"); err == nil || err.(*fs.PathError).Err != ErrNotEmptyDir {
		t.Fatal("should bee ErrNotEmptyDir error")
	}

	// try to remove /what which has children entry
	if err := memFs.Remove("/what"); err != nil{
		t.Fatal("should no error here")
	}
}