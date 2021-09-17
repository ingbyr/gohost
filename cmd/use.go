/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	useSimulateFlag bool

	useCmd = &cobra.Command{
		Use:   "use",
		Short: "use group host as system host",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			host.M.ApplyGroup(args[0], useSimulateFlag)
		},
	}
)

func init() {
	useCmd.PersistentFlags().BoolVarP(&useSimulateFlag, "simulate", "s", false, "just print host content")
}