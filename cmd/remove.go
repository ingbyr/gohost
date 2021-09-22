/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/conf"
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
	"strings"
)

var (
	removeGroup bool
	removeCmd   = &cobra.Command{
		Use:   "rm",
		Short: "Delete host or group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if removeGroup {
				host.M.DeleteGroups(strings.Split(args[0], conf.SepInCmd))
			} else {
				host.M.DeleteHosts(strings.Split(args[0], conf.SepInCmd))
			}
		},
	}
)

func init() {
	removeCmd.Flags().BoolVarP(&removeGroup, "group", "g", false, "remove group")
}
