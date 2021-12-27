/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	useSimulateFlag bool

	// TODO use multi groups
	useCmd = &cobra.Command{
		Use:   "use",
		Short: "use group host as system host",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			app.ApplyGroup(args[0], useSimulateFlag)
		},
	}
)

func init() {
	useCmd.PersistentFlags().BoolVarP(&useSimulateFlag, "simulate", "s", false, "just print host content")
}
