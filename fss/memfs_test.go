/*
 @Author: ingbyr
*/

package fss

import (
	"github.com/google/go-cmp/cmp"
	"io/fs"
	"testing"
)

func TestAssertHostFs(t *testing.T) {
	var _ HostFs = NewMemFs()
}

func TestMemFsUtil_ValidatePath(t *testing.T) {
	var tests = []struct{
		path string
		want string
	} {
		{"", ""},
		{"./", ""},
		{"..", ""},
		{"/", "mem"},
		{"/a/bb/ccc/d","mem/a/bb/ccc/d"},
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

func TestMemFs_ReadDir(t *testing.T) {
	var memFs HostFs = NewMemFs()
	dirs := []string {
		"/d1/d11",
		"/d1/d12/d111",
		"/d2",
	}
	for _, dir := range dirs {
		if err := memFs.MkdirAll(dir, fs.ModeDir | 0644); err != nil {
			t.Fatal(err)
		}
	}
	if dirs, err := memFs.ReadDir("/d1"); err != nil {
		t.Fatal(err)
	} else {
		for _, dir := range dirs {
			t.Log(dir.Name(), dir.IsDir())
		}
	}
}

func TestMemFs(t *testing.T) {
	var memFs HostFs = NewMemFs()
	dir := "/x/y"
	if err := memFs.MkdirAll(dir, fs.ModeDir | 0644); err != nil {
		t.Fatal(err)
	}
	filePath := dir + "/name.txt"
	fileContent := []byte("from ingbyr")
	if err := memFs.WriteFile(filePath, fileContent, 0664); err != nil {
		t.Fatal(err)
	}
	file, err := memFs.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	t.Log(file)

	stat, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(stat)

	var result = make([]byte, int(stat.Size()))
	n, err := file.Read(result)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(result[:n],fileContent); diff != "" {
		t.Error("diff", diff)
	}

	if dirs, err := memFs.ReadDir("/"); err != nil {
		t.Fatal(err)
	} else {
		for _, dir := range dirs {
			t.Log(dir.Name(), dir.IsDir())
		}
	}

}
