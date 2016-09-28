package cmd

import (
	"github.com/revolvingcow/pair/shell"
	"github.com/spf13/cobra"
)

var routerCmd = &cobra.Command{
	Use:   "router",
	Short: "An SSH router for pairing sessions",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		shell.Router(hostKey, localPort)
	},
}
