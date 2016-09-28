package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

func HostKey(path string) (ssh.Signer, error) {
	expanded := expand(path)

	// Try to read in the PEM private key generated for this application
	raw, err := ioutil.ReadFile(expanded)
	if err != nil {
		// If the private key does not exit the we try to generate and
		// save it
		private, err := rsa.GenerateKey(rand.Reader, 2048)
		file, err := os.Create(expanded)
		if err != nil {
			return nil, fmt.Errorf("Failed to create private key")
		}

		// Create the PEM block
		pk := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(private),
		}

		// Encode so we may read it back in on further runs
		err = pem.Encode(file, pk)
		file.Close()

		// Read the file back in for good measure
		raw, err = ioutil.ReadFile(expanded)
		if err != nil {
			return nil, fmt.Errorf("Failed to read generated private key")
		}
	}

	key, err := ssh.ParsePrivateKey(raw)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key")
	}

	return key, nil
}
