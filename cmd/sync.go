package cmd

import (
	"fmt"

	"github.com/revolvingcow/pair/keys"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync [username...]",
	Short: "Synchronize user(s) public keys",
	Long:  `Synchronize user(s) authorized public keys from popular external sources (Github and Gitlab).`,
	Run: func(cmd *cobra.Command, args []string) {
		store := keys.NewKeystore()
		store.Sync(command, args...)
		store.Save()
		fmt.Println(store.String("* ", args...))
	},
}
