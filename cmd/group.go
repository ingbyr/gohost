/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
	"strings"
)

var (
	gdel string
	gadd string

	groupCommand = &cobra.Command{
		Use:   "group",
		Short: "manage host group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hostName := args[0]
			if gdel != "" {
				host.Manager.DeleteGroup(hostName, strings.Split(gdel, ","))
			}
			if gadd != "" {
				host.Manager.AddGroup(hostName, strings.Split(gadd, ","))
			}
			host.Manager.PrintGroup(hostName)
		},
	}
)

func init() {
	groupCommand.Flags().StringVarP(&gdel, "delete", "d", "", "delete some groups from host")
	groupCommand.Flags().StringVarP(&gadd, "add", "a", "", "add some groups to host")
}
