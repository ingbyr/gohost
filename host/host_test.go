/*
 @Author: ingbyr
*/

package host

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ingbyr/gohost/myfs"
	"testing"
)

func TestNewHostByFileName(t *testing.T) {
	var tests = []struct {
		fileName string
		want     *Host
	}{
		{"", &Host{Name: "", FileName: "", Groups: []string{}}},
		{"dev", &Host{Name: "dev", FileName: "dev", Groups: []string{}}},
		{"dev_dev", &Host{Name: "dev", FileName: "dev_dev", Groups: []string{"dev"}}},
		{"dev_test_prod_host1", &Host{Name: "host1", FileName: "dev_test_prod_host1", Groups: []string{"dev", "test", "prod"}}},
		{"dev_test_prod_test_host1", &Host{Name: "host1", FileName: "dev_test_prod_test_host1", Groups: []string{"dev", "test", "prod", "test"}}},
	}
	for _, test := range tests {
		host := NewHostByFileName(test.fileName)
		if diff := cmp.Diff(test.want, host, cmpopts.IgnoreFields(Host{}, "FilePath")); diff != "" {
			t.Logf("input %v\n", test.fileName)
			t.Logf("diff \n%s\n", diff)
			t.Fail()
		}
	}
}

func TestManager_CreateRemoveNewHost(t *testing.T) {
	M.SetFs(myfs.NewMemFs())
	M.CreateNewHost("dev1", []string{"dev"}, false)
	M.LoadHosts()
	M.PrintGroups()
	M.PrintHosts()
}
