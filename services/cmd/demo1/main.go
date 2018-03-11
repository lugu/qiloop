package main

import (
	"encoding/json"
	"fmt"
	"github.com/lugu/qiloop/services"
	"github.com/lugu/qiloop/session/dummy"
	"log"
)

func Print(i interface{}) {
	json, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		log.Fatalf("json encoding failed: %s", err)
	}
	fmt.Println(string(json))
}

func main() {
	sess, err := dummy.NewSession(":9559")
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	logManager, err := services.NewLogManager(sess, 1)
	if err != nil {
		log.Fatalf("failed to connect log manager: %s", err)
	}
	objRef, err := logManager.CreateListener()
	if err != nil {
		log.Fatalf("failed to get listener: %s", err)
	}
	obj, err := sess.Object(objRef)
	if err != nil {
		log.Fatalf("failed to connect listener: %s", err)
	}
	meta, err := obj.MetaObject(objRef.ObjectID)
	if err != nil {
		log.Fatalf("failed to register event: %s", err)
	}
	Print(meta)
}
