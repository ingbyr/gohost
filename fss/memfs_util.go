/*
 @Author: ingbyr
*/

package fss

import "io/fs"

// validPath path must start with root "/" and not contains any ".", ".."
// return "mem/path"
func validPath(path string) string {
	invalidPath := ""
	rootPath := "mem"
	if path == "" || path == "." || path == ".." {
		return invalidPath
	}
	if path == "/" {
		return rootPath
	}
	path = rootPath + path
	if fs.ValidPath(path) {
		return path
	}
	return invalidPath
}
