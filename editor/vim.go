package editor

import (
	"fmt"
	"os"
	"os/exec"
)

func OpenByVim(filePath string) error {
	const vim = "vim"
	const noSwap = "-n"
	if _, err := exec.LookPath(vim); err != nil {
		return fmt.Errorf("please install vim before editing file\n")
	}
	cmd := exec.Command(vim, noSwap, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
