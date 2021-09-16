/*
 @Author: ingbyr
*/

package fss

import "unicode/utf8"

// validPath input must start with root "/" and not contains any ".", ".."
func validPath(path string) bool {
	if path == "" || !utf8.ValidString(path) || path[0] != '/'{
		return false
	}
	if path == "/" {
		return true
	}
	path = path[1:]
	for {
		i := 0
		for i < len(path) && path[i] != '/' {
			i++
		}
		part := path[:i]
		if part == "" || part == "." || part == ".." {
			return false
		}
		if i == len(path) {
			return true
		}
		path = path[i+1:]
	}
}
