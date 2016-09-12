package main

import (
	"fmt"
	"os"

	"github.com/revolvingcow/pair/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if cmd.GenerateDocumentation {
		if err := doc.GenMarkdownTree(cmd.RootCmd, "./doc"); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}
