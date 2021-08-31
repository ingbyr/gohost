/*
 @Author: ingbyr
*/

package cmd

import (
	"fmt"
	"github.com/ingbyr/gohost/host"
	"github.com/spf13/cobra"
	"strings"
)

var (
	listCmd = &cobra.Command{
		Use:   "ls",
		Short: "List all group host",
		Run: func(cmd *cobra.Command, args []string) {
			_, groups := host.Manager.LoadHostNodes()
			for group, nodes := range groups {
				var sb strings.Builder
				for _, node := range nodes {
					sb.WriteString(node.Name)
					sb.WriteString(", ")
				}
				nodeNames := sb.String()[:sb.Len()-2]
				fmt.Printf("%s: %s", group, nodeNames)
				fmt.Printf("\n")
			}
		},
	}

	listAllCmd = &cobra.Command{
		Use:   "a",
		Short: "List all host file",
		Run: func(cmd *cobra.Command, args []string) {
			nodes, _ := host.Manager.LoadHostNodes()
			fmt.Println("All hosts: ")
			for name, node := range nodes {
				fmt.Printf("  %s (%s)\n", name, node.FileName)
			}
		},
	}
)

func init() {
	listCmd.AddCommand(listAllCmd)
}
