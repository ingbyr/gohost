/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	newCmd = &cobra.Command{
		Use:   "new",
		Short: "create new host file",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			host.Manager.CreateNewHost(args[0], args[1:])
		},
	}
)
