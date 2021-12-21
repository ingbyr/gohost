//go:build linux || darwin
// +build linux darwin

package conf

import (
	"golang.org/x/text/encoding/unicode"
)

const (
	SysHost = "/etc/hosts"
	NewLine = "\n"
)

var (
	SysHostCharset = unicode.UTF8
)
