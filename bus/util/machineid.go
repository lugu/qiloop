package util

import (
	"log"
	"os"

	"github.com/denisbrodbeck/machineid"
)

// MachineID returns a fixed ID specific to an application.
func MachineID() string {
	id, err := machineid.ProtectedID("qi-messaging")
	if err != nil {
		log.Fatalf("OS not supported: failed to read machine-id: %s", err)
	}
	return id
}

// ProcessID return the PID.
func ProcessID() uint32 {
	return uint32(os.Getpid())
}
