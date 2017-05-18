package cmd

import "github.com/spf13/cobra"

var (
	// GenerateDocumentation is a flag to determine if markdown documentation should be created.
	GenerateDocumentation bool

	// Command is a custom command to be executed upon connection
	command string

	hostKey    string
	local      string
	localPort  int
	router     string
	privateKey string
)

const (
	// AppVersion is the applications current version
	AppVersion = "0.0.3"
)

func init() {
	RootCmd.Flags().BoolVar(&GenerateDocumentation, "doc", false, "Generate documentation for all commands")

	addCmd.Flags().StringVarP(&command, "command", "c", "", "Custom command to be executed upon connection")
	syncCmd.Flags().StringVarP(&command, "command", "c", "", "Custom command to be executed upon connection")
	hostCmd.Flags().StringVar(&local, "local", "localhost", "")
	hostCmd.Flags().StringVar(&router, "mux", "mux.revolvingcow.com:443", "")
	hostCmd.Flags().StringVar(&privateKey, "private-key", "~/.ssh/id_rsa", "")
	joinCmd.Flags().StringVar(&privateKey, "private-key", "~/.ssh/id_rsa", "")
	joinCmd.Flags().StringVar(&router, "mux", "mux.revolvingcow.com:443", "")
	routerCmd.Flags().StringVar(&hostKey, "host-key", "~/.ssh/pair", "Path to the host private key")
	routerCmd.Flags().IntVar(&localPort, "port", 443, "Port to listen on. Zero (0) for random.")

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
