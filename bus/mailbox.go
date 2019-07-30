package bus

import (
	"log"

	"github.com/lugu/qiloop/bus/net"
)

// Mail contains a message to which one can respond.
type Mail struct {
	Msg  *net.Message
	From Channel
}

// NewMail returns a new Mail
func NewMail(msg *net.Message, from Channel) Mail {
	return Mail{
		Msg:  msg,
		From: from,
	}
}

// MailBox is a FIFO for messages
type MailBox chan Mail

// NewMailBox creates a mailbox and a goroutine which uses the
// Receiver to handle the incomming messages.
func NewMailBox(r Receiver) MailBox {
	box := MailBox(make(chan Mail, 10))
	go func() {
		for {
			mail, ok := <-box
			if !ok {
				return
			}
			err := r.Receive(mail.Msg, mail.From)
			log.Printf("error while processing %v: %v",
				mail.Msg.Header, err)
		}
	}()
	return box
}
