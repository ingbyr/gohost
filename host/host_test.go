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

func TestNewHostByFileName(t *testing.T) {
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
		if diff := cmp.Diff(test.want, host, cmpopts.IgnoreFields(Host{},"FilePath")); diff != "" {
			fmt.Printf("input %v\n", test.fileName)
			fmt.Printf("diff \n%s\n", diff)
			t.Fail()
		}
	}
}

func TestManager_CreateNewHost(t *testing.T) {

}