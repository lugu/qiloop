# Notes on qiloop implementation

## EndPoint abstraction

The package bus/net offers an abstraction of an established network connection
called EndPoint. The EndPoint interface allow multiple parties to send and
receive messages.

    type EndPoint interface {
    	Send(m Message) error
    	MakeHandler(f Filter, q chan<-Message, cl Closer) int
    	RemoveHandler(id int) error
    	[...]
    }

Sending a message is done by writing to a channel. In order to receive
a message one needs to provide an handler composed of 3 elements:

-   a `Filter` method to select the messages.
-   a `Queue` to receive incomming messages.
-   a `Closer` method to be notified when the connection closes.

The messages selected by Filter are sent to the queue. Writing to the
queue shall not block otherwise the message is dropped. Messages are
filtered and process in the respective order of their arrival.

## Client proxy handler

Call data flow: message is sent to the endpoint. reply message is read
by endpoint, sent to the consumer queue, extracted from the queue and
sent to the reply channel (client.go). Then read from the reply
channel, deserialized by the specialized proxy and returned to the
caller.

Subscribe data flow: event message read by endpoint, sent to the
consumer queue, extracted from the queue and sent to the event channel
(client.go). Then read from the event channel, deserialized and sent to
the subscriber channel. Then read from the subscriber's channel.

## Server message dispatch

Each incoming connection generates a new endpoint. Traffic from the various
endpoints is process in parallel.

Incoming messages from an endpoint are sent to the router (see
Router.Receive). The router dispatches the messages to the various
services. Each service dispatch the messages the mailbox of the
object. Each object has its mailbox which is process sequentially.

At this stage the messages have been sent to handler queue of the
endpoint and then to the object mailbox.

This has two consequences:

    - messages from an endpoint are post in order to the mailbox.
    - objects receive one message at a time.
