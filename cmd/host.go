package cmd

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/gooops/easyssh"
	"github.com/revolvingcow/pair/shell"
	"github.com/spf13/cobra"
)

var hostCmd = &cobra.Command{
	Use:   "as [url to Git service]/[username]",
	Short: "Declares a host for a pairing session",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("No room found")
		}

		// Start a daemon
		daemon := make(chan int)
		go shell.Daemon(hostKey, daemon)

		select {
		case localPort = <-daemon:
			log.Printf("Local port used: %d", localPort)
		}

		// Get the remote port and use it in the reverse tunnel
		room := strings.TrimSpace(args[0])
		if strings.Index(router, ":") == -1 {
			router = fmt.Sprintf("%s:22", router)
		}

		config := &ssh.ClientConfig{
			User: room,
			Auth: []ssh.AuthMethod{
				ssh.Password(fmt.Sprintf("%d", localPort)),
			},
		}

		client, err := easyssh.Dial("tcp", router, config)
		if err != nil {
			log.Fatalln(err)
		}
		defer client.Close()
		log.Println("Connection closed")
	},
}
