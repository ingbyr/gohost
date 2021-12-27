/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	maxLine int
	sysCmd  = &cobra.Command{
		Use:   "sys",
		Short: "display current system host",
		Run: func(cmd *cobra.Command, args []string) {
			app.PrintSysHost(maxLine)
		},
	}
)

func init() {
	sysCmd.Flags().IntVarP(&maxLine, "line", "l", 0, "display first x lines")
}
