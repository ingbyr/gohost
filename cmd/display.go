/*
 @Author: ingbyr
*/

package cmd

import "fmt"

func printError(err error)  {
	if err != nil {
		fmt.Printf("[error] %s\n", err.Error())
	}
}
