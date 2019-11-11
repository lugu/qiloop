package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aybabtme/rgbterm"
	qilog "github.com/lugu/qiloop/bus/logger"
	"github.com/lugu/qiloop/bus/session"
)

func label(l qilog.LogLevel) (color, label string) {
	switch l {
	case qilog.LogLevelFatal:
		return "{#0000ff}", "[F]"
	case qilog.LogLevelError:
		return "{#ff0000}", "[E]"
	case qilog.LogLevelWarning:
		return "{#ff8800}", "[W]"
	case qilog.LogLevelInfo:
		return "{#ffcc00}", "[I]"
	case qilog.LogLevelVerbose:
		return "{#bbbbbb}", "[V]"
	case qilog.LogLevelDebug:
		return "{#ffffff}", "[D]"
	default:
		return "{#ff0000}", "[?]"
	}
}

func printColor(m qilog.LogMessage) {
}

func logger(serverURL string, level uint32) {

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)

	stat, err := os.Stdout.Stat()
	if err != nil {
		log.Fatal(err)
	}
	mode := stat.Mode()

	colored := true
	if (mode&os.ModeDevice == 0) || (mode&os.ModeCharDevice == 0) {
		colored = false
	}

	sess, err := session.NewSession(serverURL)
	if err != nil {
		log.Fatalf("connect: %s", err)
	}
	defer sess.Terminate()
	srv := qilog.Services(sess)
	closer := func(err error) {
		log.Printf("Remote connection closed: %s", err)
		close(signalChannel)

	}
	logManager, err := srv.LogManager(closer)
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

	for {
		select {
		case _ = <-signalChannel:
			return
		case m := <-logs:
			if m.Level == qilog.LogLevelNone {
				return
			}
			color, info := label(m.Level)
			nocolor := "{}"
			out := rgbterm.ColorOut
			if !colored {
				color = ""
				nocolor = ""
				out = os.Stdout
			}
			sec := int64(m.SystemDate.Ns) / int64(time.Second)
			ns := int64(m.SystemDate.Ns) - sec*int64(time.Second)
			t := time.Unix(sec, ns).Format("2006/01/02 15:04:05.000")
			fmt.Fprintf(out, "%s%s %s %d %s %s%s\n",
				color, info, t,
				m.Id, m.Category, m.Message, nocolor)
		}
	}
}
