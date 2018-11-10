package util

import (
	"github.com/denisbrodbeck/machineid"
	"log"
	"os"
)

func MachineID() string {
	id, err := machineid.ProtectedID("qi-messaging")
	if err != nil {
		log.Fatalf("OS not supported: failed to read machine-id: %s", err)
	}
	return id
}

func ProcessID() uint32 {
	return uint32(os.Getpid())
}
