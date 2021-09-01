/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	addCmd = &cobra.Command{
		Use: "add",
		Short: "Add host file",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			host.Manager.AddHost(args[0], args[1:])
		},
	}
)
