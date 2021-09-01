/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	_group bool
	_host  bool
	mvCmd  = &cobra.Command{
		Use:   "mv",
		Short: "change host file name or groups",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if _group {
				host.Manager.ChangeGroups(args[0], args[1:])
			} else if _host {
				host.Manager.ChangeHostName(args[0], args[1])
			}
		},
	}
)

func init() {
	mvCmd.Flags().BoolVarP(&_host, "host", "t", true, "change host name")
	mvCmd.Flags().BoolVarP(&_group, "group", "g", false, "change the groups of host")
}
