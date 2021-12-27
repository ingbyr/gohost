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
	newCmd = &cobra.Command{
		Use:   "new",
		Short: "create new host file",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				app.CreateNewHost(args[0], []string{}, true)
			} else {
				app.CreateNewHost(args[0], strings.Split(args[1], config.SepInCmd), true)
			}
		},
	}
)
