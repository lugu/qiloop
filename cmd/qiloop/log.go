package main

import (
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func logger(serverURL string) {

	sess, err := session.NewSession(serverURL)
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
	srv := services.Services(sess)
	logManager, err := srv.LogManager()
	if err != nil {
		log.Fatalf("failed to access LogManager service: %s", err)
	}
	logListener, err := logManager.CreateListener()
	if err != nil {
		log.Fatalf("failed to create listener: %s", err)
	}
	defer logListener.Terminate(logListener.ObjectID())

	err = logListener.ClearFilters()
	if err != nil {
		log.Fatalf("failed to clear filters: %s", err)
	}
	cancel, logs, err := logListener.SubscribeOnLogMessage()
	if err != nil {
		log.Fatalf("failed to subscribe logs: %s", err)
	}
	defer cancel()

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)

	for {
		select {
		case _ = <-signalChannel:
			return
		case log := <-logs:
			Print(log)
		}
	}
}
