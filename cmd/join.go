package cmd

import (
	"log"
	"strings"

	"github.com/revolvingcow/pair/shell"
	"github.com/spf13/cobra"
)

var joinCmd = &cobra.Command{
	Use:   "with [url to Git service]/[username]",
	Short: "Join a hosted pairing session",
	Long:  `Join a registered pairing session over an SSH connection`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("No room found")
		}

		port := "22"
		if idx := strings.Index(router, ":"); idx != -1 {
			port = router[idx+1:]
			router = router[:idx]
		}

		shell.Client(args[0], privateKey, router, port)
	},
}
