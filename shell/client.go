package shell

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/gooops/easyssh"
	"golang.org/x/crypto/ssh"
)

func Client(room, privateKey, address, port string) {
	forward := ""
	username := "pair"
	user, err := user.Current()
	if err == nil {
		username = user.Username
	}

	config := &ssh.ClientConfig{
		User: room,
		Auth: []ssh.AuthMethod{
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
				if len(questions) > 0 {
					forward = questions[0]
				}

				return []string{"10-4"}, nil
			}),
		},
	}

	// Connect to the remote server to get the forwarding address
	log.Println("Looking for pairing session...")
	conn, err := easyssh.Dial("tcp", fmt.Sprintf("%s:%s", address, port), config)
	if err == nil {
		conn.Close()
	}

	if forward == "" {
		log.Fatal("No hosted session for pairing found")
	} else {
		// forward = strings.Replace(forward, "localhost", address, 1)
		// forward = strings.Replace(forward, "[::]", address, 1)

		forward = forward[strings.LastIndex(forward, ":")+1:]
	}

	cmd := exec.Command("ssh", "-p", forward, fmt.Sprintf("%s@%s", username, address))
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		log.Println(err)
	}

	// // Redeclare our connection settings
	// config = &ssh.ClientConfig{
	// 	User: username,
	// 	Auth: []ssh.AuthMethod{},
	// }

	// // Check to see if we can use the private key
	// if authKey := keys.PublicKeyFile(privateKey); authKey != nil {
	// 	log.Println("Enabling use of private key", privateKey)
	// 	config.Auth = append(config.Auth, authKey)
	// }

	// // Check if they have a secure agent
	// if authAgent := keys.SSHAgent(); authAgent != nil {
	// 	log.Println("Enabling use of SSH agent")
	// 	config.Auth = append(config.Auth, authAgent)
	// }

	// // Connect to the remote server
	// log.Println("Joining session for", room)
	// conn, err = easyssh.Dial("tcp", forward, config)
	// if err != nil {
	// 	log.Fatalf("Unable to connect: %s", err)
	// }
	// defer conn.Close()

	// // Resize terminal so it displays things nicely
	// fd := int(os.Stdin.Fd())
	// oldState, err := terminal.MakeRaw(fd)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// // Create a new SSH session
	// session, err := conn.NewSession()
	// if err != nil {
	// 	log.Fatalf("Could not create new session: %s", err)
	// }

	// // Some finalization stuff when we return to our local shell
	// finalize := func() {
	// 	session.Close()
	// 	terminal.Restore(fd, oldState)
	// }
	// defer finalize()

	// // Set up I/O
	// session.Stderr = os.Stderr
	// session.Stdout = os.Stdout
	// session.Stdin = os.Stdin

	// // Get the terminal height and width
	// termWidth, termHeight, err := terminal.GetSize(fd)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// // Some terminal mode configuration
	// modes := ssh.TerminalModes{
	// 	ssh.ECHO:          1,
	// 	ssh.TTY_OP_ISPEED: 14400,
	// 	ssh.TTY_OP_OSPEED: 14400,
	// }

	// // Request a pseudo terminal
	// if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
	// 	log.Fatalf("Request for a pseudo terminal failed: %s", err)
	// }

	// // Launch a login shell
	// if err := session.Shell(); err != nil {
	// 	log.Fatalf("Could not start shell: %s", err)
	// }

	// // Wait for an exit from the shell
	// _ = session.Wait()
	log.Println("Connection closed")
}
