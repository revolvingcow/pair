package cmd

import (
	"fmt"

	"github.com/revolvingcow/pair/keys"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [username...]",
	Short: "Remove user(s) public keys",
	Long:  `Remove all, or some, authorized public access keys associated with the user account`,
	Run: func(cmd *cobra.Command, args []string) {
		store := keys.NewKeystore()
		fmt.Println(store.String("- ", args...))
		store.Remove(args...)
		store.Save()
	},
}
