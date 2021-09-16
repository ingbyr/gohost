/*
 @Author: ingbyr
*/

package fss

import (
	"github.com/google/go-cmp/cmp"
	"io/fs"
	"testing"
)

func TestMemFsUtil_ValidatePath(t *testing.T) {
	var tests = []struct{
		input string
		want  bool
	} {
		{"", false},
		{"./", false},
		{"..", false},
		{"/", true},
		{"/a/bb/ccc/d", true},
		{"/a/bbb/ccc/", false},
		{"/a/bb/ccc/./d", false},
		{"/a/bbb/../ccc", false},
	}

	for _, test := range tests {
		if diff := cmp.Diff(test.want, validPath(test.input)); diff != "" {
			t.Fatalf("\n input %s\n diff: %s", test.input, diff)
		}
	}
}

func TestMemFs(t *testing.T) {
	memFs := NewMemFs()
	dir := "/x/y"
	if err := memFs.MkdirAll(dir, fs.ModeDir | 0644); err != nil {
		t.Fatal(err)
	}
	filePath := "x/y/name.txt"
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


	t.Log(memFs.ReadDir(""))
}
