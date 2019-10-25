package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aybabtme/rgbterm"
	qilog "github.com/lugu/qiloop/bus/logger"
	"github.com/lugu/qiloop/bus/session"
)

func level(l qilog.LogLevel) string {
	switch l {
	case qilog.LogLevelFatal:
		return "{#ff0000}[FATAL]"
	case qilog.LogLevelError:
		return "{#ff0000}[ERROR]"
	case qilog.LogLevelWarning:
		return "{#ffaa00}[WARN ]"
	case qilog.LogLevelInfo:
		return "{#00ff00}[INFO ]"
	case qilog.LogLevelVerbose:
		return "{#ffffff}[VERB ]"
	case qilog.LogLevelDebug:
		return "{#0000ff}[DEBUG]"
	default:
		return "{#ff0000}[UNEXP]"
	}
}

func printLog(m qilog.LogMessage) {
	if m.Level == qilog.LogLevelNone {
		return
	}
	fmt.Fprintln(rgbterm.ColorOut, level(m.Level),
		m.Category, m.Source, m.Message, "{}")
}

func logger(serverURL string, level uint32) {

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

	err = logListener.ClearFilters()
	if err != nil {
		log.Fatalf("clear filters: %s", err)
	}
	cancel, logs, err := logListener.SubscribeOnLogMessage()
	if err != nil {
		log.Fatalf("subscribe logs: %s", err)
	}
	defer cancel()

	err = logListener.SetLevel(qilog.LogLevel{int32(level)})
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
			printLog(log)
		}
	}
}
