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
				app.AddGroup(hostName, strings.Split(groupAdd, config.SepInCmd))
			}
			if groupDel != "" {
				app.DeleteHostGroups(hostName, strings.Split(groupDel, config.SepInCmd))
			}
			if groupList {
				app.PrintGroup(hostName)
			}
		},
	}
)

func init() {
	groupCommand.Flags().BoolVarP(&groupList, "list", "l", true, "list host groups")
	groupCommand.Flags().StringVarP(&groupDel, "delete", "d", "", "delete some groups from host")
	groupCommand.Flags().StringVarP(&groupAdd, "add", "a", "", "add some groups to host")
}
