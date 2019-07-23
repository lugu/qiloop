package logger

import (
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
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

		message(LogLevelNone, "qi.none", "boom"),
		message(LogLevelFatal, "qi.fatal", "boom"),
		message(LogLevelError, "qi.error", "boom"),
		message(LogLevelInfo, "qi.info", "boom"),
		message(LogLevelWarning, "qi.warning", "boom"),
		message(LogLevelVerbose, "qi.verbose", "boom"),
		message(LogLevelDebug, "qi.debug", "boom"),

		message(LogLevelNone, "qiloop.none", "boom"),
		message(LogLevelFatal, "qiloop.fatal", "boom"),
		message(LogLevelError, "qiloop.error", "boom"),
		message(LogLevelInfo, "qiloop.info", "boom"),
		message(LogLevelWarning, "qiloop.warning", "boom"),
		message(LogLevelVerbose, "qiloop.verbose", "boom"),
		message(LogLevelDebug, "qiloop.debug", "boom"),
	}
}

func TestLogListener(t *testing.T) {

	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	ns := bus.PrivateNamespace()
	srv, err := bus.StandAloneServer(listener, bus.Yes{}, ns)
	if err != nil {
		t.Fatal(err)
	}

	session := srv.Session()
	_, err = srv.NewService("LogManager", NewLogManager())
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Terminate()

	proxies := Services(session)
	logManager, err := proxies.LogManager()
	if err != nil {
		t.Fatal(err)
	}

	logListener, err := logManager.CreateListener()
	if err != nil {
		t.Fatal(err)
	}

	err = logListener.SetVerbosity(LogLevelWarning)
	if err != nil {
		t.Fatal(err)
	}
	err = logListener.SetCategory("qi.*", LogLevelVerbose)
	if err != nil {
		t.Fatal(err)
	}
	err = logListener.SetFilters(map[string]LogLevel{
		"qiloop.": LogLevelInfo,
		"boo":     LogLevelDebug,
	})
	if err != nil {
		t.Fatal(err)
	}
	cancel, messages, err := logListener.SubscribeOnLogMessage()
	if err != nil {
		t.Fatal(err)
	}

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
			}
			if strings.Contains(msg.Category, "qi.") {
				if msg.Level.Level > LogLevelVerbose.Level {
					t.Errorf("filter not respected: %#v", msg)
				}
			} else if strings.Contains(msg.Category, "qiloop.") {
				if msg.Level.Level > LogLevelInfo.Level {
					t.Errorf("filter not respected: %#v", msg)
				}
			} else if msg.Level.Level > LogLevelWarning.Level {
				t.Errorf("default level not respected: %#v", msg)
			}
		}
		wait.Done()
	}()
	cancel()
	wait.Wait()

	err = logListener.ClearFilters()
	if err != nil {
		t.Error(err)
	}
}

func TestLogProvider(t *testing.T) {

	addr := util.NewUnixAddr()
	listener, err := net.Listen(addr)
	if err != nil {
		t.Fatal(err)
	}
	ns := bus.PrivateNamespace()
	srv, err := bus.StandAloneServer(listener, bus.Yes{}, ns)
	if err != nil {
		t.Fatal(err)
	}

	session := srv.Session()
	_, err = srv.NewService("LogManager", NewLogManager())
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Terminate()

	constructor := Services(session)
	logManager, err := constructor.LogManager()
	if err != nil {
		t.Fatal(err)
	}

	service := logManager.ProxyService(session)
	impl, logger := newLogProviderImpl("cat")
	providerProxy, err := constructor.NewLogProvider(service, impl)
	id, err := logManager.AddProvider(providerProxy)
	if err != nil {
		t.Fatal(err)
	}
	defer logManager.RemoveProvider(id)

	logger2, err := NewLogger(session, "test")

	logListener, err := logManager.GetListener()
	if err != nil {
		t.Fatal(err)
	}
	err = logListener.SetVerbosity(LogLevelVerbose)
	if err != nil {
		t.Fatal(err)
	}

	cancel, messages, err := logListener.SubscribeOnLogMessage()
	if err != nil {
		t.Fatal(err)
	}

	var messageCount int
	var wait sync.WaitGroup

	wait.Add(1)
	go func() {
		for msg := range messages {
			if msg.Level.Level == LogLevelNone.Level {
				t.Errorf("none level message: %#v", msg)
			} else if msg.Level.Level > LogLevelVerbose.Level {
				t.Errorf("default level not respected: %#v", msg)
			}
			messageCount++
		}
		wait.Done()
	}()

	logger.Error("paf")
	logger.Warning("paf")
	logger.Info("paf")
	logger.Verbose("paf")
	logger2.Debug("paf")
	logger2.Debug("pif")
	logger2.Debug("pof")

	cancel()
	wait.Wait()

	if messageCount != 4 {
		t.Errorf("wrong number of messages (%d != 4)", messageCount)
	}
}
