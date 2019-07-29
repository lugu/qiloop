package token

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

// AuthFile contains the user credentials. The format is:
// First line: user login
// Second line: user token
// Leading and trailing spaces are trimmed.
var AuthFile = ".qiloop-auth.conf"

func init() {
	usr, err := user.Current()
	if err == nil {
		AuthFile = usr.HomeDir + string(os.PathSeparator) + AuthFile
	}
}

// GetUserToken returns user login and token.
func GetUserToken() (string, string) {
	return readUserToken()
}

func readUserToken() (string, string) {
	file, err := os.Open(AuthFile)
	if err != nil {
		return "", ""
	}
	defer file.Close()
	r := bufio.NewReader(file)
	user, err := r.ReadString('\n')
	if err != nil {
		return "", ""
	}
	pwd, err := r.ReadString('\n')
	if err != io.EOF && err != nil {
		return "", ""
	}
	return strings.TrimSpace(user), strings.TrimSpace(pwd)
}

// WriteUserToken save the user credentials.
func WriteUserToken(login string, token string) error {

	var flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(AuthFile, flag, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open auth file: %s", err)
	}
	_, err = file.WriteString(login + "\n" + token + "\n")
	if err != nil {
		return fmt.Errorf("Failed to write auth file: %s", err)
	}

	return nil
}
