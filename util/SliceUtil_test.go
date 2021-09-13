/*
 @Author: ingbyr
*/

package util

import (
	"fmt"
	"testing"
)

func TestUnion(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []string{"c", "d"}
	fmt.Println(SliceUnion(s1, s2))
	fmt.Println(SliceUnion(s2, s1))
}
