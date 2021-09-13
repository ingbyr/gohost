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

func Err(err error) {
	if err != nil {
		fmt.Printf("[error] %s\n", err.Error())
	}
}

func ErrExit(err error) {
	if err != nil {
		fmt.Printf("[error] %s\n", err.Error())
		os.Exit(1)
	}
}

func Panic(msg string, err error) {
	fmt.Printf("[panic] %s\n", msg)
	ErrExit(err)
}
