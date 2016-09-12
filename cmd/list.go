package cmd

import (
	"fmt"

	"github.com/revolvingcow/pair/keys"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [username...]",
	Short: "List user(s) and their authorized public keys",
	Long:  `List will display all, or some, of the public keys and the associated user.`,
	Run: func(cmd *cobra.Command, args []string) {
		store := keys.NewKeystore()
		fmt.Println(store.String("", args...))
	},
}
