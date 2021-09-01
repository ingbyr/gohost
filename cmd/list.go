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
		Short: "glist all group",
		Run: func(cmd *cobra.Command, args []string) {
			host.Manager.PrintGroups()
		},
	}

	listAllCmd = &cobra.Command{
		Use:   "all",
		Short: "glist all host file",
		Run: func(cmd *cobra.Command, args []string) {
			host.Manager.PrintHosts()
		},
	}
)

func init() {
	listCmd.AddCommand(listAllCmd)
}
