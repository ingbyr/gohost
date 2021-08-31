/*
 @Author: ingbyr
*/

package host

import (
	"github.com/ingbyr/gohost/display"
	"strings"
)

type Line struct {
	Domain string
	Ip     string
}

func parseLine(line string) *Line {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	fields := strings.Fields(line)
	if len(fields) < 2 {
		display.Warn("skip invalid line: " + line)
		return nil
	}
	return &Line{
		Domain: strings.TrimSpace(fields[0]),
		Ip:     strings.TrimSpace(fields[1]),
	}
}
