package net

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/lugu/qiloop/basic"
	"io"
)

// Magic is a constant to discriminate between message and garbage.
const Magic uint32 = 0x42dead42

// Version is the supported version of the protocol.
const Version uint16 = 0

// Message types:
const (
	Unknown uint8 = iota
	Call
	Reply
	Error
	Post
	Event
	Capability
	Cancel
	Cancelled
)

// HeaderSize is the size of a message header. It is the
// minimum size of a message.
const HeaderSize uint32 = 28

// Header represents a message header.
type Header struct {
	Magic   uint32 // magic number
	ID      uint32 // message id
	Size    uint32 // size of the payload
	Version uint16 // protocol version
	Type    uint8  // type of the message
	Flags   uint8  // flags
	Service uint32 // service id
	Object  uint32 // object id
	Action  uint32 // function or event id
}

// NewHeader construct a message header given some parameters. The size
// of the message is zero.
func NewHeader(typ uint8, service uint32, object uint32, action uint32, id uint32) Header {
	return Header{
		Magic, id, 0, Version, typ, 0, service, object, action,
	}
}

func (h *Header) writeMagic(w io.Writer) error {
	buf := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(buf, h.Magic)
	if bytes, err := w.Write(buf); err != nil {
		return err
	} else if bytes != 4 {
		return fmt.Errorf("failed to write magic (%d instead of 4)", bytes)
	}
	return nil
}

func (h *Header) Write(w io.Writer) (err error) {
	wrap := func(field string, err error) error {
		return fmt.Errorf("failed to write message %s: %s", field, err)
	}
	if err = h.writeMagic(w); err != nil {
		return wrap("magic", err)
	}
	if err = basic.WriteUint32(h.ID, w); err != nil {
		return wrap("id", err)
	}
	if err = basic.WriteUint32(h.Size, w); err != nil {
		return wrap("size", err)
	}
	if err = basic.WriteUint16(h.Version, w); err != nil {
		return wrap("version", err)
	}
	if err = basic.WriteUint8(h.Type, w); err != nil {
		return wrap("type", err)
	}
	if err = basic.WriteUint8(h.Flags, w); err != nil {
		return wrap("flags", err)
	}
	if err = basic.WriteUint32(h.Service, w); err != nil {
		return wrap("service", err)
	}
	if err = basic.WriteUint32(h.Object, w); err != nil {
		return wrap("object", err)
	}
	if err = basic.WriteUint32(h.Action, w); err != nil {
		return wrap("action", err)
	}
	return nil
}

func (h *Header) readMagic(r io.Reader) error {
	buf := []byte{0, 0, 0, 0}
	if bytes, err := r.Read(buf); err != nil {
		return err
	} else if bytes != 4 {
		return fmt.Errorf("failed to read magic (%d instead of 4)", bytes)
	} else {
		h.Magic = binary.BigEndian.Uint32(buf)
		return nil
	}
}

// Read parses a message header from an io.Reader.
func (h *Header) Read(r io.Reader) (err error) {
	if err = h.readMagic(r); err != nil {
		return fmt.Errorf("failed to read message magic: %s", err)
	} else if h.Magic != Magic {
		return fmt.Errorf("invalid message magic: %x", h.Magic)
	}
	if h.ID, err = basic.ReadUint32(r); err != nil {
		return fmt.Errorf("failed to read message id: %s", err)
	}
	if h.Size, err = basic.ReadUint32(r); err != nil {
		return fmt.Errorf("failed to read message size: %s", err)
	}
	if h.Version, err = basic.ReadUint16(r); err != nil {
		return fmt.Errorf("failed to read message version: %s", err)
	} else if h.Version != Version {
		return fmt.Errorf("invalid message version: %d", h.Version)
	}
	if h.Type, err = basic.ReadUint8(r); err != nil {
		return fmt.Errorf("failed to read message type: %s", err)
	} else if h.Type == Unknown || h.Type > Cancelled {
		return fmt.Errorf("invalid message type: %d", h.Type)
	}
	if h.Flags, err = basic.ReadUint8(r); err != nil {
		return fmt.Errorf("failed to read message flags: %s", err)
	}
	if h.Service, err = basic.ReadUint32(r); err != nil {
		return fmt.Errorf("failed to read message service: %s", err)
	}
	if h.Object, err = basic.ReadUint32(r); err != nil {
		return fmt.Errorf("failed to read message object: %s", err)
	}
	if h.Action, err = basic.ReadUint32(r); err != nil {
		return fmt.Errorf("failed to read message action: %s", err)
	}
	return nil
}

// Message represents a QiMessaging message.
type Message struct {
	Header  Header
	Payload []byte
}

// Write marshal a message into an io.Writer. The header and the
// payload are written in a single write operation.
func (m *Message) Write(w io.Writer) error {

	if uint32(len(m.Payload)) != m.Header.Size {
		return fmt.Errorf("invalid message size: %d instead of %d", len(m.Payload), m.Header.Size)
	}

	// Pack header and payload in a buffer and then it to the network.
	buf := bytes.NewBuffer(make([]byte, 0, HeaderSize+m.Header.Size))

	if err := m.Header.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize header: %s", err)
	}

	if size, err := buf.Write(m.Payload); err != nil {
		return fmt.Errorf("failed to write message payload: %s", err)
	} else if size != int(m.Header.Size) {
		return fmt.Errorf("failed to write message payload (%d instead of %d)", size, m.Header.Size)
	}

	if size, err := w.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to send message: %s", err)
	} else if size != int(m.Header.Size+HeaderSize) {
		return fmt.Errorf("message not completly wrote (%d instead of %d)", size, m.Header.Size+HeaderSize)
	}
	return nil
}

// Read unmarshal a message from io.Reader. First the header is read,
// then if correct the payload is read. The payload will not be read
// if the header is not considerred well formatted.
func (m *Message) Read(r io.Reader) error {

	// Read the complete header, then parse the fields.
	b := make([]byte, HeaderSize)
	if size, err := r.Read(b); err != nil {
		return err // won't process reader issues.
	} else if size != int(HeaderSize) {
		// TODO: truncated message: read again.
		return fmt.Errorf("full header not received (%d instead of %d)", size, HeaderSize)
	}

	if err := m.Header.Read(bytes.NewBuffer(b)); err != nil {
		return fmt.Errorf("failed to read message header: %s", err)
	}
	// Filter messages larger than 10 MB
	if m.Header.Size > 10*1024*1024 {
		return fmt.Errorf("won't process message this large: %d", m.Header.Size)
	}
	m.Payload = make([]byte, m.Header.Size)
	size, err := r.Read(m.Payload)
	if err != nil {
		return fmt.Errorf("failed to read message payload: %s", err)
	} else if size != int(m.Header.Size) {
		return fmt.Errorf("failed to read message payload (%d instead of %d)", size, m.Header.Size)
	}
	return nil
}

// NewMessage assemble an header and a payload to create a message.
// The size filed of the header is adjusted if necessary.
func NewMessage(header Header, payload []byte) Message {
	header.Size = uint32(len(payload))
	return Message{header, payload}
}
