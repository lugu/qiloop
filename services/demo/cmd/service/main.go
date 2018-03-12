package main

import (
	"fmt"
	"github.com/lugu/qiloop/services"
	"github.com/lugu/qiloop/session/dummy"
	"github.com/lugu/qiloop/session/util"
	"log"
)

func main() {
	sess, err := dummy.NewSession(":9559")
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	directory, err := services.NewServiceDirectory(sess, 1)
	if err != nil {
		log.Fatalf("directory creation failed: %s", err)
	}
	info := services.ServiceInfo{
		Name:      "My own service",
		ServiceId: 9999,
		MachineId: util.MachineID(),
		ProcessId: util.ProcessID(),
		Endpoints: nil,
		SessionId: "", // FIXME: what is it?
	}
	serviceID, err := directory.RegisterService(info)
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}
	fmt.Print(serviceID)
}
