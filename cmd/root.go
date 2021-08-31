/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gohost",
	Short: "Host Switcher written in go",
}

func Execute() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
