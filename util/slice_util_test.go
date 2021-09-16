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
		if diff := cmp.Diff(union, test.wantUnion); diff != "" {
			fmt.Printf("input %v, %v\n", test.input[0], test.input[1])
			fmt.Printf("diff union \n%s\n", diff)
			t.Fail()
		}
		if diff := cmp.Diff(add, test.wantAdd); diff != "" {
			fmt.Printf("input %v, %v\n", test.input[0], test.input[1])
			fmt.Printf("diff add \n%s\n", diff)
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
		if diff := cmp.Diff(res, test.want); diff != "" {
			fmt.Printf("input %v\n", test.input)
			fmt.Printf("diff %s\n", diff)
			t.Fail()
		}
	}
}
