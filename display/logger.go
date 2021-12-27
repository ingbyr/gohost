/*
 @Author: ingbyr
*/

package display

import (
	"fmt"
	"os"
)

func Warn(warn string) {
	fmt.Printf("[warn] %s\n", warn)
}

func Err(errors ...error) {
	for _, err := range errors {
		fmt.Printf("[error] %s\n", err.Error())
	}
}

func ErrExit(errors ...error) {
	Err(errors...)
	os.Exit(1)
}

func ErrStrExit(msg string, errors ...error) {
	fmt.Printf("[error] %s\n", msg)
	ErrExit(errors...)
}
