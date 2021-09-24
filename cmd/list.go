/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	listAll bool
	listCmd = &cobra.Command{
		Use:   "ls",
		Short: "list all group",
		Run: func(cmd *cobra.Command, args []string) {
			if listAll {
				host.M.DisplayHosts()
			} else {
				host.M.DisplayGroups()
			}
		},
	}
)

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "List all host file")
}
