/*
 @Author: ingbyr
*/

package host

import (
	"fmt"
	"testing"
)

func TestUnion(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []string{"c", "d"}
	fmt.Println(union(s1, s2))
	fmt.Println(union(s2, s1))
}
