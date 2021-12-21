package editor

import (
	"fmt"
)

type None struct {
	name string
}

func NewNone(name string) *None {
	return &None{name: name}
}

func (none *None) Open(filePath string) error {
	return fmt.Errorf("not valid editor %s", none.name)
}
