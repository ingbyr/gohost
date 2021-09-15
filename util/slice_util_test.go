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
	s1 := []string{"a", "b", "c"}
	s2 := []string{"c", "d"}
	fmt.Println(SliceUnion(s1, s2))
	fmt.Println(SliceUnion(s2, s1))
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
			fmt.Printf("%+v\n", test)
			fmt.Printf("got %+v\n", res)
			t.Fail()
		}
	}
}
