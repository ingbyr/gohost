/*
 @Author: ingbyr
*/

package fss

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ingbyr/gohost/util"
	"io/fs"
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

func TestMemFs_CreateDirFile(t *testing.T) {
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

}

func TestMemFs(t *testing.T) {
	var memFs HostFs = NewMemFs()
	dir := "/x/y"
	if err := memFs.MkdirAll(dir, fs.ModeDir|0644); err != nil {
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

	stat, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	var result = make([]byte, int(stat.Size()))
	n, err := file.Read(result)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(result[:n], fileContent); diff != "" {
		t.Error("diff", diff)
	}

}
