/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "ls",
		Short: "List all group host",
		Run: func(cmd *cobra.Command, args []string) {
			host.Manager.PrintGroups()
		},
	}

	listAllCmd = &cobra.Command{
		Use:   "all",
		Short: "List all host file",
		Run: func(cmd *cobra.Command, args []string) {
			host.Manager.PrintHostNodes()
		},
	}
)

func init() {
	listCmd.AddCommand(listAllCmd)
}
