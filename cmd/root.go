/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/display"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gohost",
	Short: "Host switcher written in go",
}

func Execute() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(renameCmd)
	rootCmd.AddCommand(sysCmd)
	rootCmd.AddCommand(groupCommand)
	rootCmd.AddCommand(delCmd)
	if err := rootCmd.Execute(); err != nil {
		display.Err(err)
		return
	}
}
