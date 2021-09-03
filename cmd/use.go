/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	applyCmd = &cobra.Command{
		Use:   "use",
		Short: "apply group host to system host",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			host.Manager.ApplyGroup(args[0])
		},
	}
)
