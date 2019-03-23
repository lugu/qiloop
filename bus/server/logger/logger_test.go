package logger

import (
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/util"
	"regexp"
	"sync"
	"testing"
)

func matchHelper(t *testing.T, pattern, source string, expected bool) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		t.Error(err)
	}
	if reg.FindString(source) == "" {
		// not found
		if expected {
			t.Errorf("unexpected missmatch (%s, %s)", pattern, source)
		}
	} else {
		if !expected {
			t.Errorf("unexpected match (%s, %s)", pattern, source)
		}
	}
}

func TestCategoryMatch(t *testing.T) {
	matchHelper(t, "test", "test", true)
	matchHelper(t, "test", "some test", true)
	matchHelper(t, "test", "test again", true)
	matchHelper(t, "some test", "test", false)
	matchHelper(t, "test again", "test", false)
	matchHelper(t, "test*", "test", true)
	matchHelper(t, "test*", "test again", true)
	matchHelper(t, "test.*", "test again", true)
	matchHelper(t, ".*test", "test", true)
	matchHelper(t, "test*", "test again", true)
	matchHelper(t, "test*", "some test", true)
	matchHelper(t, "test*", "test", true)
	matchHelper(t, "test*", "test again", true)
}

func message(level LogLevel, category, content string) LogMessage {
	return LogMessage{
		Level:    level,
		Category: category,
		Message:  content,
	}
}

func messagesList() []LogMessage {
	return []LogMessage{
		message(LogLevelNone, "none", "boom"),
		message(LogLevelFatal, "fatal", "boom"),
		message(LogLevelError, "error", "boom"),
		message(LogLevelInfo, "info", "boom"),
		message(LogLevelWarning, "warning", "boom"),
		message(LogLevelVerbose, "verbose", "boom"),
		message(LogLevelDebug, "debug", "boom"),
	}
}

func TestLogListener(t *testing.T) {

	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	ns := server.PrivateNamespace()
	srv, err := server.StandAloneServer(listener, server.Yes{}, ns)
	if err != nil {
		t.Fatal(err)
	}

	session := srv.Session()
	_, err = srv.NewService("LogManager", NewLogManager())
	if err != nil {
		t.Fatal(err)
	}

	proxies := Services(session)
	logManager, err := proxies.LogManager()
	if err != nil {
		t.Fatal(err)
	}

	logListener, err := logManager.CreateListener()
	if err != nil {
		t.Fatal(err)
	}

	cancel, messages, err := logListener.SubscribeOnLogMessage()
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	err = logManager.Log(messagesList())
	if err != nil {
		t.Error(err)
	}

	var wait sync.WaitGroup

	wait.Add(1)
	go func() {
		for msg := range messages {
			if msg.Level.Level == LogLevelNone.Level {
				t.Errorf("none level message: %#v", msg)
			} else if msg.Level.Level > LogLevelInfo.Level {
				t.Errorf("default level not respected: %#v", msg)
			}
		}
		wait.Done()
	}()
	// FIXME: logListener.OnTerminate does not seems to be called
	// when Terminate is called.

	// FIXME: does terminate informs the signal subscribers ?
	logListener.Terminate(logListener.ObjectID())
	logManager.Terminate(logManager.ObjectID())
	srv.Terminate()
	wait.Wait()

}
