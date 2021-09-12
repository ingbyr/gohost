package editor

import (
	"github.com/ingbyr/micro/v2/cmd/micro"
)

func OpenByMicro(filePath string) error {
	micro.Open(filePath)
	return nil
}
