package cmd

import (
	"fmt"

	"github.com/revolvingcow/pair/keys"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [username...]",
	Short: "Add user(s) public keys",
	Long:  `Query popular external sources (Github and Gitlab) for user public keys which can then be added to local SSH authorized keys file.`,
	Run: func(cmd *cobra.Command, args []string) {
		store := keys.NewKeystore()
		store.Get(command, args...)
		store.Save()
		fmt.Println(store.String("+ ", args...))
	},
}
