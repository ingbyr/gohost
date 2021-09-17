/*
 @Author: ingbyr
*/

package myfs

import "io/fs"

// validPath path must start with root "/" and not contains any ".", ".."
// return "mem/path"

const invalidPath = ""

func validPath(path string) string {
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
