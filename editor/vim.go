package editor

import (
	"fmt"
	"os"
	"os/exec"
)

type Vim struct {
	cmd     string
	options []string
}

func init() {
	register("vim", newVim)
}

func newVim() Editor {
	return &Vim{
		cmd:     "vim",
		options: []string{"-n"},
	}
}

func (vim *Vim) Open(filePath string) error {
	if _, err := exec.LookPath(vim.cmd); err != nil {
		return fmt.Errorf("not find editor %s", vim.cmd)
	}
	args := append(vim.options, filePath)
	cmd := exec.Command(vim.cmd, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
