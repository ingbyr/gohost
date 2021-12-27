package editor

import (
	"fmt"
	"github.com/ingbyr/gohost/config"
	"os"
	"os/exec"
	"strings"
)

type Editor interface {
	Open(filePath string) error
}

type editor struct {
	Command string
	Args    []string
}

func New(commands string) Editor {
	fields := strings.Fields(commands)
	e := &editor{}
	e.Command = fields[0]
	if len(fields) > 1 {
		e.Args = fields[1:]
	}
	return e
}

func (e *editor) Open(filePath string) error {
	if _, err := exec.LookPath(e.Command); err != nil {
		return fmt.Errorf("not find editor %s", e.Command)
	}
	var args []string
	if e.Args != nil && len(e.Args) > 0 {
		args = append(e.Args, filePath)
	} else {
		args = []string{filePath}
	}
	cmd := exec.Command(e.Command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func ExtractArgs(args string) []string {
	args = strings.TrimSpace(args)
	if len(args) == 0 {
		return []string{}
	}
	return strings.Split(args, config.SepInCmd)
}
