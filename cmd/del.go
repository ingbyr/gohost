package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	delCmd = &cobra.Command{
		Use:   "del",
		Short: "delete some hosts",
		Run: func(cmd *cobra.Command, args []string) {
			host.Manager.DeleteHosts(args)
		},
	}
)
