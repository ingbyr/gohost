/*
 @Author: ingbyr
*/

package host

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestNewHost(t *testing.T) {
	var tests = []struct {
		hostName string
		want *Host
	}{
		{"", nil },
		{"dev", nil},
		{"dev_dev", nil},
		{"dev_test_prod_host1", nil},
		{"dev_test_prod_test_host1", nil},
	}
	for _, test := range tests {
		host := NewHostByFileName(test.hostName)
		if !cmp.Equal(test.want, fmt.Sprintf("%+v", host)) {
			fmt.Printf("%+v\n", test)
			fmt.Printf("got %+v\n", host)
			t.Fail()
		}
	}
}
