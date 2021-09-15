/*
 @Author: ingbyr
*/

package util

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestUnion(t *testing.T) {
	var tests = []struct {
		input     [2][]string
		wantUnion []string
		wantAdd   []string
	}{
		{[2][]string{{"a", "b", "c"}, {"c", "d"}}, []string{"a", "b", "c", "d"}, []string{"d"}},
		{[2][]string{{}, {"c", "d"}}, []string{"c", "d"}, []string{"c", "d"}},
		{[2][]string{{"a", "b", "b"}, {"a"}}, []string{"a", "b", "b"}, []string{}},
	}

	for _, test := range tests {
		union, add := SliceUnion(test.input[0], test.input[1])
		if !cmp.Equal(union, test.wantUnion) {
			fmt.Printf("input %v, %v\n", test.input[0], test.input[1])
			fmt.Printf("diff union\n%v\n", cmp.Diff(union, test.wantUnion))
			t.Fail()
		}
		if !cmp.Equal(add, test.wantAdd) {
			fmt.Printf("input %v, %v\n", test.input[0], test.input[1])
			fmt.Printf("diff add\n%v\n", cmp.Diff(union, test.wantUnion))
			t.Fail()
		}
	}
}

func TestSortUniqueStringSlice(t *testing.T) {
	var tests = []struct {
		input []string
		want  []string
	}{
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", "a", "a"}, []string{"a"}},
		{[]string{"a", "c", "b", "d"}, []string{"a", "b", "c", "d"}},
		{[]string{"a", "c", "a", "d", "b", "d"}, []string{"a", "b", "c", "d"}},
	}

	for _, test := range tests {
		res := SortUniqueStringSlice(test.input)
		if !cmp.Equal(res, test.want) {
			fmt.Printf("input %v\n", test.input)
			fmt.Printf("diff %v\n", cmp.Diff(res, test.want))
			t.Fail()
		}
	}
}
