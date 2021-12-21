//go:build windows
// +build windows

package conf

import "golang.org/x/text/encoding/charmap"

const (
	SysHost = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	NewLine = "\r\n"
)

var (
	// see https://stackoverflow.com/questions/701882/what-is-ansi-format/702023#702023
	SysHostCharset = charmap.ISO8859_1
)
