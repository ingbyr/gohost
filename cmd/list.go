/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use: "ls",
		Short: "List all hosts",
		Run: func(cmd *cobra.Command, args []string) {
			host.Manager.LoadHostNodes()
		},
	}
)
