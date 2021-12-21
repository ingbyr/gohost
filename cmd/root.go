/*
 @Author: ingbyr
*/

package cmd

import (
	"github.com/ingbyr/gohost/display"
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gohost",
	Short: "Host switcher made by ingbyr",
}

func Execute() {
	host.M.SetCmdMode()
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(renameCmd)
	rootCmd.AddCommand(sysCmd)
	rootCmd.AddCommand(groupCommand)
	rootCmd.AddCommand(removeCmd)
	if err := rootCmd.Execute(); err != nil {
		display.ErrExit(err)
		return
	}
}
