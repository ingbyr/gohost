/*
 @Author: ingbyr
*/

package host

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"path/filepath"
	"testing"
)

func TestNewHost(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"", "&{Name: FileName: FilePath:test Groups:[]}"},
		{"dev", "&{Name:dev FileName:dev FilePath:test/dev Groups:[]}"},
		{"dev_dev", "&{Name:dev FileName:dev_dev FilePath:test/dev_dev Groups:[dev]}"},
		{"dev_test_prod_host1", "&{Name:host1 FileName:dev_test_prod_host1 FilePath:test/dev_test_prod_host1 Groups:[dev prod test]}"},
		{"dev_test_prod_test_host1", "&{Name:host1 FileName:dev_test_prod_test_host1 FilePath:test/dev_test_prod_test_host1 Groups:[dev prod test]}"},
	}
	for _, test := range tests {
		host := NewHostByFileName(test.input, filepath.Join("test", test.input))
		if !cmp.Equal(test.want, fmt.Sprintf("%+v", host)) {
			fmt.Printf("%+v\n", test)
			fmt.Printf("got %+v\n", host)
			t.Fail()
		}
	}
}
