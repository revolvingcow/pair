package cmd

import (
	"fmt"
	"log"
	"net"
	"strings"

	"dev.justinjudd.org/justin/easyssh"

	"golang.org/x/crypto/ssh"

	"github.com/revolvingcow/pair/shell"
	"github.com/spf13/cobra"
)

var hostCmd = &cobra.Command{
	Use:   "as [url to Git service]/[username]",
	Short: "Declares a host for a pairing session",
	Long:  `Start hosting a pairing session by ensuring an SSH daemon is listening and then registering with an external service`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("No room found")
		}

		// Determine if there is a local SSH daemon already running
		// on the common port
		localPort := 22
		sshdAddress, err := net.ResolveTCPAddr("tcp", "[::]:ssh")
		sshdListener, err := net.ListenTCP("tcp", sshdAddress)
		if err == nil {
			// Close our listener
			sshdListener.Close()

			// Start a daemon
			daemon := make(chan int)
			go shell.Daemon(hostKey, daemon)

			// Wait for the local port
			select {
			case localPort = <-daemon:
				log.Printf("Spawning an SSH daemon on port %d", localPort)
			}
		} else {
			log.Println("Found a local SSH daemon running")
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
