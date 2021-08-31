package editor

import (
	"os"
	"os/exec"
)

func Open(filePath string) {
	cmd := exec.Command("vim", "-n", filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
