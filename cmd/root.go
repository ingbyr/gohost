package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	root = &cobra.Command{
		Use:   "gohost",
		Short: "Simple host switcher tool made by ingbyr",
		Run: func(cmd *cobra.Command, args []string) {
			//m := tui.New()
			//p := tea.NewProgram(m)
			//if err := p.Start(); err != nil {
			//	log.Fatal(err)
			//}
		},
	}
)

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
	}
}
