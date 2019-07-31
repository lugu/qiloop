package bus

import "github.com/lugu/qiloop/bus/net"

// mail contains a message to which one can respond.
type mail struct {
	msg  *net.Message
	from Channel
}

// mailBox is a FIFO for messages
type mailBox chan mail
