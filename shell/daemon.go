package shell

import (
	"fmt"
	"log"

	"dev.justinjudd.org/justin/easyssh"

	"github.com/revolvingcow/pair/keys"

	"golang.org/x/crypto/ssh"
)

func Daemon(hostKey string, ch chan int) {
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			_, err := keys.PublicKeyIsAuthorized(key)
			if err != nil {
				log.Printf("[User %s] Connection rejected from %s", c.User(), c.RemoteAddr())
				return nil, err
			}

			log.Printf("[User %s] Connection accepted from %s", c.User(), c.RemoteAddr())
			return &ssh.Permissions{}, nil
		},
	}

	// Parse and add the host key to be leveraged by the daemon
	pk, err := keys.HostKey(hostKey)
	if err != nil {
		log.Fatalf("Failed to parse private key")
	}
	config.AddHostKey(pk)
	log.Printf("Fingerprint: %s", keys.Fingerprint(pk.PublicKey()))

	// Only allow handling of sessions to keep it simple and reduce
	// the attack vector
	mux := easyssh.NewChannelsMux()
	mux.HandleChannel(easyssh.SessionRequest, easyssh.SessionHandler())

	handler := easyssh.NewStandardSSHServerHandler()
	handler.MultipleChannelsHandler = mux

	// Attempt to start the daemon on a random port
	for attempts := 10; attempts > 0; attempts-- {
		localPort := 1024 + portRandomizer.Intn(60000)

		server := easyssh.Server{
			Addr:    fmt.Sprintf("%s:%d", "", localPort),
			Config:  config,
			Handler: handler,
		}

		ch <- localPort
		_ = server.ListenAndServe()
		break
	}
}
