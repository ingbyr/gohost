//go:build windows
// +build windows

package conf

import "golang.org/x/text/encoding/charmap"

const (
	SysHost = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	NewLine = "\r\n"
)

var (
	SysHostCharset = charmap.ISO8859_1
)
