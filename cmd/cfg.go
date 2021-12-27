package cmd

import (
	"github.com/ingbyr/gohost/conf"
	"github.com/spf13/cobra"
)

var (
	confCmd = &cobra.Command{
		Use:   "cfg",
		Short: "change config",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			conf.Change(args[0], args[1])
		},
	}
)
