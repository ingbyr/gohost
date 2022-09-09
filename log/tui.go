package log

import (
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

var logFile *os.File

func New(file string) error {
	var err error
	logFile, err = tea.LogToFile(file, "gohost")
	if err != nil {
		return err
	}
	Debug("==============")
	return nil
}

func Debug(msg string) {
	if _, err := logFile.WriteString("[Debug] " + msg + "\n"); err != nil {
		panic(err)
	}
}
