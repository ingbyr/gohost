/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/config"
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
				app.DeleteGroups(strings.Split(args[0], config.SepInCmd))
			} else {
				app.DeleteHosts(strings.Split(args[0], config.SepInCmd))
			}
		},
	}
)

func init() {
	removeCmd.Flags().BoolVarP(&removeGroup, "group", "g", false, "remove group")
}
