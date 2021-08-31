/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	editCmd = &cobra.Command{
		Use: "edit",
		Short: "Edit one host file",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := host.Manager.EditHostFile(args[0])
			printError(err)
		},
	}
)
