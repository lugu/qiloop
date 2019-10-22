package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	qilog "github.com/lugu/qiloop/bus/logger"
	"github.com/lugu/qiloop/bus/session"
)

func logger(serverURL string) {

	sess, err := session.NewSession(serverURL)
	if err != nil {
		log.Fatalf("connect: %s", err)
	}
	srv := qilog.Services(sess)
	logManager, err := srv.LogManager()
	if err != nil {
		log.Fatalf("access LogManager service: %s", err)
	}
	logListener, err := logManager.CreateListener()
	if err != nil {
		log.Fatalf("create listener: %s", err)
	}
	defer logListener.Terminate(logListener.ObjectID())

	m, err := logListener.MetaObject(0)
	fmt.Printf("Meta: %s\n", m.JSON())

	err = logListener.ClearFilters()
	if err != nil {
		log.Fatalf("clear filters: %s", err)
	}
	cancel, logs, err := logListener.SubscribeOnLogMessage()
	if err != nil {
		log.Fatalf("subscribe logs: %s", err)
	}
	defer cancel()

	err = logListener.SetLevel(qilog.LogLevelDebug)
	if err != nil {
		log.Fatalf("set verbosity: %s", err)
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)

	for {
		select {
		case _ = <-signalChannel:
			return
		case log := <-logs:
			fmt.Printf("%s\n", log.String())
		}
	}
}
