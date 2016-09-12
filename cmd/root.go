package cmd

import "github.com/spf13/cobra"

var (
	// GenerateDocumentation is a flag to determine if markdown documentation should be created.
	GenerateDocumentation bool

	// Command is a custom command to be executed upon connection
	Command string
)

func init() {
	RootCmd.Flags().BoolVar(&GenerateDocumentation, "doc", false, "Generate documentation for all commands")
	RootCmd.PersistentFlags().StringVarP(&Command, "command", "c", "", "Custom command to be executed upon connection")
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(removeCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(syncCmd)
}

// RootCmd is the entry point for the application from which all actions are subcommands.
var RootCmd = &cobra.Command{
	Use:   "pair",
	Short: "Pair aims to simplify the pair programming experience",
	Long:  `Pair aims to simplify the pair programming experience.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
