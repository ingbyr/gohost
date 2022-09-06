package main

import (
	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use:   "gohost",
		Short: "Simple gohost switcher tool made by ingbyr",
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

}
