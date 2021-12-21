package editor

import (
	"fmt"
	"os"
	"os/exec"
)

type vim struct{}

func NewVim() vim {
	return vim{}
}

func (vim) Open(filePath string) error {
	const vimCmd = "vim"
	const noSwap = "-n"
	if _, err := exec.LookPath(vimCmd); err != nil {
		return fmt.Errorf("can not find vim")
	}
	cmd := exec.Command(vimCmd, noSwap, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
