/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/conf"
	"github.com/ingbyr/gohost/host"
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
				host.Manager.CreateNewHost(args[0], []string{})
			} else {
				host.Manager.CreateNewHost(args[0], strings.Split(args[1], conf.SepGroupInCmd))
			}
		},
	}
)
