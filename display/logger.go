/*
 @Author: ingbyr
*/

package display

import (
	"fmt"
	"os"
)

func Warn(warn string) {
	if warn != "" {
		fmt.Printf("[warn] %s\n", warn)
	}
}

func Err(errors ...error) {
	for _, err := range errors {
		fmt.Printf("[error] %s\n", err.Error())
	}
}

func ErrExit(errors ...error) {
	for _, err := range errors {
		fmt.Printf("[error] %s\n", err.Error())
	}
	os.Exit(1)
}

func Panic(msg string, err error) {
	fmt.Printf("[panic] %s\n", msg)
	ErrExit(err)
}
