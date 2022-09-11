package log

import (
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"time"
)

var logFile *os.File

func New(file string) error {
	var err error
	logFile, err = tea.LogToFile(file, "gohost")
	if err != nil {
		return err
	}
	Debug(">>> Start logging at " + nowStr())
	return nil
}

func Debug(msg string) {
	if _, err := logFile.WriteString("[DEBUG " + nowStr() + "] " + msg + "\n"); err != nil {
		panic(err)
	}
}

func Error(err error) {

	if _, err := logFile.WriteString("[ERROR" + nowStr() + "] " + err.Error()); err != nil {
		panic(err)
	}
}

func nowStr() string {
	return time.Now().Format(time.RFC3339)
}
