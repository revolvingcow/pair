package cmd

import "github.com/spf13/cobra"

var (
	// GenerateDocumentation is a flag to determine if markdown documentation should be created.
	GenerateDocumentation bool

	// Command is a custom command to be executed upon connection
	Command string

	hostKey    string
	local      string
	localPort  int
	router     string
	privateKey string
)

const (
	// AppVersion is the applications current version
	AppVersion = "0.0.1"
)

func init() {
	RootCmd.Flags().BoolVar(&GenerateDocumentation, "doc", false, "Generate documentation for all commands")

	addCmd.Flags().StringVarP(&Command, "command", "c", "", "Custom command to be executed upon connection")
	syncCmd.Flags().StringVarP(&Command, "command", "c", "", "Custom command to be executed upon connection")
	hostCmd.Flags().StringVar(&local, "local", "localhost", "")
	hostCmd.Flags().StringVar(&router, "router", "mux.revolvingcow.com:443", "")
	hostCmd.Flags().StringVar(&privateKey, "private-key", "~/.ssh/id_rsa", "")
	joinCmd.Flags().StringVar(&privateKey, "private-key", "~/.ssh/id_rsa", "")
	joinCmd.Flags().StringVar(&router, "router", "mux.revolvingcow.com:443", "")
	routerCmd.Flags().StringVar(&hostKey, "host-key", "~/.ssh/pair", "Path to the host private key")
	routerCmd.Flags().IntVar(&localPort, "port", 443, "Port for the router to listen on. Zero (0) for random.")

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(removeCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(syncCmd)
	RootCmd.AddCommand(hostCmd)
	RootCmd.AddCommand(joinCmd)
	RootCmd.AddCommand(routerCmd)
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
