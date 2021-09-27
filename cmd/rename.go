/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	renameGroupFlag bool
	renameCmd = &cobra.Command{
		Use:   "mv",
		Short: "rename host file name",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if renameGroupFlag {
				host.M.ChangeGroupName(args[0], args[1])
			} else {
				host.M.ChangeHostName(args[0], args[1])
			}
		},
	}
)

func init() {
	renameCmd.Flags().BoolVarP(&renameGroupFlag, "group", "g", false, "rename group name")
}
