package shell

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"dev.justinjudd.org/justin/easyssh"

	"github.com/revolvingcow/pair/keys"
	"github.com/revolvingcow/pair/store"

	"golang.org/x/crypto/ssh"
)

func Router(hostKey string, localPort int) {
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			return setupProxy(c, pass, hostKey)
		},
		KeyboardInteractiveCallback: func(c ssh.ConnMetadata, client ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {
			room := &store.Room{Name: c.User()}
			if err := room.Get(); err == nil {
				_, err = client(c.User(), "Please redirect your connection to one of the following addresses", []string{room.LocalAddress}, []bool{false})
			}

			return nil, fmt.Errorf("Redirect request sent")
		},
	}

	pk, err := keys.HostKey(hostKey)
	if err != nil {
		log.Fatalf("Failed to parse private key")
	}
	config.AddHostKey(pk)
	log.Printf("Fingerprint: %s", keys.Fingerprint(pk.PublicKey()))

	easyssh.HandleChannel(easyssh.SessionRequest, easyssh.SessionHandler())
	easyssh.HandleChannel(easyssh.DirectForwardRequest, easyssh.DirectPortForwardHandler())
	easyssh.HandleRequestFunc(easyssh.RemoteForwardRequest, easyssh.TCPIPForwardRequest)

	if localPort != 0 {
		if err = easyssh.ListenAndServe(fmt.Sprintf("%s:%d", "[::]", localPort), config, nil); err != nil {
			log.Fatalf("Error: %s", err)
		}
		return
	}

	for attempts := 10; attempts > 0; attempts-- {
		port := 1024 + portRandomizer.Intn(60000)
		err = easyssh.ListenAndServe(fmt.Sprintf("%s:%d", "[::]", port), config, nil)
		if err != nil {
			log.Printf("Error: %s", err)
			continue
		}

		break
	}
}

type proxyLogger struct {
	room *store.Room
}

func (l *proxyLogger) Trace(f string, args ...interface{}) {
	// log.Printf(f, args...)
}

func (l *proxyLogger) Debug(f string, args ...interface{}) {
	// log.Printf(f, args...)
}

func (l *proxyLogger) Info(f string, args ...interface{}) {
	if strings.HasPrefix(f, "Closed (") {
		l.room.Connections--
		l.room.Update()
	} else if strings.HasPrefix(f, "Opened") {
		l.room.Connections++
		l.room.Update()
	}
}

func (l *proxyLogger) Warn(f string, args ...interface{}) {
	if strings.HasPrefix(f, "Remote connection failed:") {
		l.room.Eject()
	}
}

func setupProxy(c ssh.ConnMetadata, pass []byte, hostKey string) (*ssh.Permissions, error) {
	user := c.User()
	room := &store.Room{Name: user}
	if err := room.Get(); err != nil || room == nil {
		room.Name = user
		room.Connections = 0
		room.Created = time.Now()
		room.Touched = time.Now()
	}

	for attempts := 10; attempts > 0; attempts-- {
		port := 1024 + portRandomizer.Intn(60000)
		room.LocalAddress = fmt.Sprintf("[::]:%d", port)

		localAddress, err := net.ResolveTCPAddr("tcp", room.LocalAddress)
		if err != nil {
			log.Println("Error connecting to local address:", err)
			continue
		}

		room.RemoteAddress = c.RemoteAddr().String()
		idx := strings.LastIndex(room.RemoteAddress, ":")
		if idx == -1 {
			room.RemoteAddress = fmt.Sprintf("%s:%s", room.RemoteAddress, string(pass))
		} else {
			room.RemoteAddress = fmt.Sprintf("%s:%s", room.RemoteAddress[:idx], string(pass))
		}

		remoteAddress, err := net.ResolveTCPAddr("tcp", room.RemoteAddress)
		if err != nil {
			log.Println("Error connecting to remote address:", err)
			break
		}

		listener, err := net.ListenTCP("tcp", localAddress)
		if err != nil {
			log.Println("Error listening on:", err)
			break
		}

		room.Update()

		log.Printf("Listening on %s for %s", room.LocalAddress, room.RemoteAddress)
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Println("Warning listening on:", err)
				continue
			}

			conn2, err := net.DialTCP("tcp", nil, remoteAddress)
			if err != nil {
				log.Println("Warning dialing remote:", err)
				continue
			}

			go proxy(conn2, conn)
		}
	}

	return nil, nil
}

func proxy(remote, local *net.TCPConn) {
	remoteClosed := make(chan struct{}, 1)
	localClosed := make(chan struct{}, 1)

	go broker(remote, local, localClosed)
	go broker(local, remote, remoteClosed)

	// var waitFor chan struct{}
	// select {
	// case <-localClosed:
	// 	remote.SetLinger(8)
	// 	remote.CloseRead()
	// 	waitFor = remoteClosed

	// case <-remoteClosed:
	// 	local.CloseRead()
	// 	waitFor = localClosed
	// }

	// <-waitFor
}

func broker(dst, src net.Conn, srcClosed chan struct{}) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Printf("Copy error: %s", err)
	}

	if err := src.Close(); err != nil {
		log.Printf("Close error: %s", err)
	}

	srcClosed <- struct{}{}
}
