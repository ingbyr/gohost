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
	groupList bool
	groupDel  string
	groupAdd  string

	groupCommand = &cobra.Command{
		Use:   "cg",
		Short: "change host group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hostName := args[0]
			if groupAdd != "" {
				host.Manager.AddGroup(hostName, strings.Split(groupAdd, conf.SepInCmd))
			}
			if groupDel != "" {
				host.Manager.DeleteHostGroups(hostName, strings.Split(groupDel, conf.SepInCmd))
			}
			if groupList {
				host.Manager.PrintGroup(hostName)
			}
		},
	}
)

func init() {
	groupCommand.Flags().BoolVarP(&groupList, "list", "l", true, "list host groups")
	groupCommand.Flags().StringVarP(&groupDel, "delete", "d", "", "delete some groups from host")
	groupCommand.Flags().StringVarP(&groupAdd, "add", "a", "", "add some groups to host")
}
