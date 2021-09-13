/*
 @Author: ingbyr
*/

package display

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func Table(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.Render()
}
