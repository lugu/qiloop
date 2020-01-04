package util

import (
	"os"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/denisbrodbeck/machineid"
)

var id string

func init() {
	var err error
	id, err = machineid.ProtectedID("qi-messaging")
	if err != nil {
		noise := make([]byte, 64)
		rand.Read(noise)
		mac := hmac.New(sha256.New, noise)
		id = hex.EncodeToString(mac.Sum(nil))
	}
}

// MachineID returns a fixed ID specific to an application.
func MachineID() string {
	return id
}

// ProcessID return the PID.
func ProcessID() uint32 {
	return uint32(os.Getpid())
}
