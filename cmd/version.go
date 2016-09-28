package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the current version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s\n", AppVersion)
	},
}
