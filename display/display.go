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
		os.Exit(1)
	}
}
