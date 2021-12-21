/*
 @Author: ingbyr
*/

package host

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
)

func TestNewHostByFileName(t *testing.T) {
	M.SetMockMode()
	var tests = []struct {
		name string
		want *Host
	}{
		{".txt", &Host{Name: "", FileName: ".txt", Groups: []string{}}},
		{"dev.txt", &Host{Name: "dev", FileName: "dev.txt", Groups: []string{}}},
		{"dev_dev.txt", &Host{Name: "dev", FileName: "dev_dev.txt", Groups: []string{"dev"}}},
		{"dev_test_prod_host1.txt", &Host{Name: "host1", FileName: "dev_test_prod_host1.txt", Groups: []string{"dev", "test", "prod"}}},
		{"dev_test_prod_test_host1.txt", &Host{Name: "host1", FileName: "dev_test_prod_test_host1.txt", Groups: []string{"dev", "test", "prod", "test"}}},
	}
	for _, test := range tests {
		host := NewHostByFileName(test.name)
		if diff := cmp.Diff(test.want, host, cmpopts.IgnoreFields(Host{}, "FilePath")); diff != "" {
			t.Logf("input %v\n", test.name)
			t.Logf("diff \n%s\n", diff)
			t.Fail()
		}
	}
}
