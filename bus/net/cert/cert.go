package cert

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/user"
	"strings"
)

func GenerateCertificate() (tls.Certificate, error) {

	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template,
		&template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	return tls.X509KeyPair(certPEM, keyPEM)
}

func Certificate() (cert tls.Certificate, err error) {
	certFile, keyFile, err := getCertFiles()
	if err != nil {
		return cert, err
	}
	return tls.LoadX509KeyPair(certFile, keyFile)
}

// getCertFiles returns the public and an X509 certificate and the
// associated private RSA key
func getCertFiles() (string, string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", "", fmt.Errorf("cannot get home directory")
	}
	file, err := os.Open(usr.HomeDir + "/.qi-cert.conf")
	if err != nil {
		return "", "", fmt.Errorf("cannot open $HOME/.qi-cert.conf")
	}
	defer file.Close()
	r := bufio.NewReader(file)
	cert, err := r.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("cannot read cert from $HOME/.qi-cert.conf")
	}
	key, err := r.ReadString('\n')
	if err != io.EOF && err != nil {
		return "", "", fmt.Errorf("cannot read key from $HOME/.qi-cert.conf")
	}
	return strings.TrimSpace(cert), strings.TrimSpace(key), nil
}
