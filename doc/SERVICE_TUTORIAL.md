# How to create a new service

This guide will show you how to create and test a new service.

## Define the service interface

The public interface of the service is described by an IDL file.

In this tutorial, we will implements a timestamp service.

Create a file clock.qi.idl with the content:

        package clock

        interface Timestamp
            fn nanoseconds() -> int64
        end

## Generate the server stub

The IDL file is used to generate the necessary boiler code to
implement the service. Use `qiloop` to generate the server stub:

        qiloop stub --idl clock.qi.idl --output clock_stub_gen.go

## Automate the stub generation

In order to easily update the clock_stub_gen.go file, create a file
called generate.go with the following content:

        //go:generate qiloop stub --idl clock.qi.idl --output clock_stub_gen.go
        package clock

Then execute the command:

        go generate

## Implement the service

The file clock_stub_gen.go defines an interface called
TimestampImplementor which describe the methods needed to create the
service:

        type TimestampImplementor interface {
                Activate(activation bus.Activation, helper TimestampSignalHelper) error
                OnTerminate()
                Nanoseconds() (int64, error)
        }

The `Activate` method is called just before the service registration.
In contains runtime information useful for the service. For example
the `activation` parameter contains a session to connect other
services.

The `OnTerminate` method is called just before the service
terminates.

Create a file called clock.go with the following implementation:

        package clock

        import (
                "time"

                bus "github.com/lugu/qiloop/bus"
        )

        type timestampService struct{}

        func (t timestampService) Activate(activation bus.Activation,
                helper TimestampSignalHelper) error {
                return nil
        }
        func (t timestampService) OnTerminate() {
        }
        func (t timestampService) Nanoseconds() (int64, error) {
                return time.Now().UnixNano(), nil
        }

## Create a program

In order to use the timestamp service, we need a program which uses
`NewTimestampObject` and registers it. The clock example contains such
program: see [main.go](https://github.com/lugu/qiloop/blob/master/examples/clock/cmd/service).

In order to test it, we need a running instance of QiMessaging. We can
create one with the `qiloop directory` command:

        $ qiloop directory
        2019/07/15 22:57:09 Listening at tcp://localhost:9559

Now we can start the timestamp service with:

        $ go run ./examples/clock/cmd/service/main.go
	2019/07/15 23:00:20 Service Timestamp registered

Let see if the timestamp service is advertized by the service directory using
`qiloop info`:

	$ qiloop info
	[
	    {
		"Name": "ServiceDirectory",
		"ServiceId": 1,
		"MachineId": "e9b7594a1f209b898e7a3caea5e3199a407cf5bb08d090419e4fffdeddcf167f",
		"ProcessId": 17179,
		"Endpoints": [
		    "tcp://localhost:9559"
		],
		"SessionId": ""
	    },
	    {
		"Name": "Timestamp",
		"ServiceId": 2,
		"MachineId": "e9b7594a1f209b898e7a3caea5e3199a407cf5bb08d090419e4fffdeddcf167f",
		"ProcessId": 18596,
		"Endpoints": [
		    "unix:///tmp/qiloop-271149288"
		],
		"SessionId": ""
	    }
	]
