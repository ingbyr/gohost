/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	// TODO rename group name
	renameCmd = &cobra.Command{
		Use:   "mv",
		Short: "rename host file name",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			host.M.ChangeHostName(args[0], args[1])
		},
	}
)
