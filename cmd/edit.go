/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/display"
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	editCmd = &cobra.Command{
		Use:   "edit",
		Short: "edit one host file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := host.Manager.EditHostFile(args[0])
			display.Err(err)
		},
	}
)
