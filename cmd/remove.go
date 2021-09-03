/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	removeGroup bool
	removeCmd   = &cobra.Command{
		Use:   "rm",
		Short: "Delete host or group",
		Run: func(cmd *cobra.Command, args []string) {
			if removeGroup {
				host.Manager.DeleteGroups(args)
			} else {
				host.Manager.DeleteHostsByNames(args)
			}
		},
	}
)

func init() {
	removeCmd.Flags().BoolVarP(&removeGroup, "group", "g", false, "remove group")
}
