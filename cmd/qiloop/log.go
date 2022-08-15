package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	qilog "github.com/lugu/qiloop/bus/logger"
	"github.com/lugu/qiloop/bus/session"
)

func label(l qilog.LogLevel) (label string) {
	switch l {
	case qilog.LogLevelFatal:
		return "[F]"
	case qilog.LogLevelError:
		return "[E]"
	case qilog.LogLevelWarning:
		return "[W]"
	case qilog.LogLevelInfo:
		return "[I]"
	case qilog.LogLevelVerbose:
		return "[V]"
	case qilog.LogLevelDebug:
		return "[D]"
	default:
		return "[?]"
	}
}

func printColor(m qilog.LogMessage) {
}

func logger(serverURL string, level uint32) {

	sess, err := session.NewSession(serverURL)
	if err != nil {
		log.Fatalf("connect: %s", err)
	}
	defer sess.Terminate()
	logManager, err := qilog.LogManager(sess)
	if err != nil {
		log.Fatalf("access LogManager service: %s", err)
	}
	logListener, err := logManager.CreateListener()
	if err != nil {
		log.Fatalf("create listener: %s", err)
	}
	defer logListener.Terminate(logListener.Proxy().ObjectID())

	err = logListener.ClearFilters()
	if err != nil {
		log.Fatalf("clear filters: %s", err)
	}
	cancel, logs, err := logListener.SubscribeOnLogMessages()
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
		case msgs := <-logs:
			for _, m := range msgs {
				if m.Level == qilog.LogLevelNone {
					return
				}
				info := label(m.Level)
				sec := int64(m.SystemDate.Ns) / int64(time.Second)
				ns := int64(m.SystemDate.Ns) - sec*int64(time.Second)
				t := time.Unix(sec, ns).Format("2006/01/02 15:04:05.000")
				fmt.Printf("%s %s %d %s %s\n",
					info, t, m.Id, m.Category, m.Message)
			}
		}
	}
}
