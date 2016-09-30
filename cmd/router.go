package cmd

import (
	"github.com/revolvingcow/pair/shell"
	"github.com/spf13/cobra"
)

var routerCmd = &cobra.Command{
	Use:   "router",
	Short: "An SSH router for pairing sessions",
	Long:  `Run an SSH router to allow pairing hosts to register sessions and to connect participating parties`,
	Run: func(cmd *cobra.Command, args []string) {
		shell.Router(hostKey, localPort)
	},
}
