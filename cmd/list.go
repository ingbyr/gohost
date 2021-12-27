/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	listAll bool
	listCmd = &cobra.Command{
		Use:   "ls",
		Short: "list all group",
		Run: func(cmd *cobra.Command, args []string) {
			if listAll {
				app.DisplayHosts()
			} else {
				app.DisplayGroups()
			}
		},
	}
)

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "List all host file")
}
