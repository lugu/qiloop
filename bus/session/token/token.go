package token

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

var userLogin = ""
var userToken = ""

func init() {
	userLogin, userToken = readUserToken()
}

// GetUserToken returns user login and token.
func GetUserToken() (string, string) {
	return userLogin, userToken
}

func readUserToken() (string, string) {
	usr, err := user.Current()
	if err != nil {
		return "", ""
	}
	file, err := os.Open(usr.HomeDir + "/.qi-auth.conf")
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
	// FIXME: don't trim space from pwd
	return strings.TrimSpace(user), strings.TrimSpace(pwd)
}

// WriteUserToken save the user credentials.
func WriteUserToken(login string, token string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	var flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(usr.HomeDir+"/.qi-auth.conf", flag, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open auth file: %s", err)
	}
	_, err = file.WriteString(login + "\n" + token + "\n")
	if err != nil {
		return fmt.Errorf("Failed to write auth file: %s", err)
	}

	userLogin = login
	userToken = token

	return nil
}
