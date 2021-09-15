/*
 @Author: ingbyr
*/

package host

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
)

func TestNewHost(t *testing.T) {
	var tests = []struct {
		fileName string
		want     *Host
	}{
		{"", &Host{Name: "", FileName: "", Groups: []string{}} },
		{"dev", &Host{Name: "dev", FileName: "dev", Groups: []string{}}},
		{"dev_dev", &Host{Name: "dev", FileName: "dev_dev", Groups: []string{"dev"}}},
		{"dev_test_prod_host1", &Host{Name: "host1", FileName: "dev_test_prod_host1", Groups: []string{"dev", "test", "prod"}}},
		{"dev_test_prod_test_host1", &Host{Name: "host1", FileName: "dev_test_prod_test_host1", Groups: []string{"dev", "test", "prod", "test"}}},
	}
	for _, test := range tests {
		host := NewHostByFileName(test.fileName)
		if !cmp.Equal(test.want, host, cmpopts.IgnoreFields(Host{},"FilePath")) {
			fmt.Printf("input %v\n", test.fileName)
			fmt.Printf("diff\n%+v\n", cmp.Diff(test.want, host, cmpopts.IgnoreFields(Host{},"FilePath")))
			t.Fail()
		}
	}
}
