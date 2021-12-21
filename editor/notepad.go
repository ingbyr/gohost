//go:build windows
// +build windows

package editor

import (
	"fmt"
	"os/exec"
)

type Notepad struct {
	cmd string
}

func init() {
	register("notepad", newNotepad)
}
func newNotepad() Editor {
	return &Notepad{
		cmd: "notepad",
	}
}

func (notepad *Notepad) Open(filePath string) error {
	if _, err := exec.LookPath(notepad.cmd); err != nil {
		return fmt.Errorf("not found editor %s", notepad.cmd)
	}
	cmd := exec.Command(notepad.cmd, filePath)
	return cmd.Run()
}
