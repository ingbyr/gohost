/*
 @Author: ingbyr
*/

package display

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
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

func Table(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.Render()
}
