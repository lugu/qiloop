package cert

import (
	"bufio"
	"io"
	"os"
	"os/user"
	"strings"
)

// GetCertKey returns the public and an X509 certificate and the
// associated private RSA key
func GetCertKey() (string, string) {
	usr, err := user.Current()
	if err != nil {
		return "", ""
	}
	file, err := os.Open(usr.HomeDir + "/.qi-cert.conf")
	if err != nil {
		return "", ""
	}
	defer file.Close()
	r := bufio.NewReader(file)
	cert, err := r.ReadString('\n')
	if err != nil {
		return "", ""
	}
	key, err := r.ReadString('\n')
	if err != io.EOF && err != nil {
		return "", ""
	}
	return strings.TrimSpace(cert), strings.TrimSpace(key)
}
