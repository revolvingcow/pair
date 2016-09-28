package keys

import (
	"crypto/md5"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

func Fingerprint(k ssh.PublicKey) string {
	bytes := md5.Sum(k.Marshal())
	strbytes := make([]string, len(bytes))
	for i, b := range bytes {
		strbytes[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(strbytes, ":")
}
