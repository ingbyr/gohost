/*
 @Author: ingbyr
*/

package display

import "fmt"

func Warn(warn string)  {
	if warn != "" {
		fmt.Printf("[warn] %s\n", warn)
	}
}

func Err(err error)  {
	if err != nil {
		fmt.Printf("[error] %s\n", err.Error())
	}
}


