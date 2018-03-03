package message

import (
    "fmt"
    "io"
    "metaqi/basic"
)

const Magic uint32 = 0x42dead42
const Version uint16 = 0

const (
    // message types are:
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

type Header struct {
    Magic       uint32 // magic number
    Id          uint32 // message id
    Size        uint32 // size of the payload
    Version     uint16 // protocol version
    Type        uint8  // type of the message
    Flags       uint8  // flags
    Service     uint32 // service id
    Object      uint32 // object id
    Action      uint32 // function or event id
}

func NewHeader(typ uint8, service uint32, object uint32, action uint32, id uint32) (h Header) {
    h.Magic = Magic
    h.Type = typ
    h.Service = service
    h.Object = object
    h.Action = action
    h.Id = id
    h.Flags = 0
    h.Size = 0
    h.Version = Version
    return
}

func (h Header) Write(w io.Writer) (err error) {
    wrap := func (field string, err error) error {
        return fmt.Errorf("failed to write message %s: %s", field, err)
    }
    if err = basic.WriteUint32(h.Magic, w); err != nil {
        return wrap("magic", err)
    }
    if err = basic.WriteUint32(h.Id, w); err != nil {
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

func (h Header) Read(r io.Reader) (err error) {
    if h.Magic, err = basic.ReadUint32(r); err != nil {
        return fmt.Errorf("failed to read message magic: %s", err)
    } else if h.Magic != Magic {
        return fmt.Errorf("invalid message magic: %s", h.Magic)
    }
    if h.Id, err = basic.ReadUint32(r); err != nil {
        return fmt.Errorf("failed to read message id: %s", err)
    }
    if h.Size, err = basic.ReadUint32(r); err != nil {
        return fmt.Errorf("failed to read message size: %s", err)
    }
    if h.Version, err = basic.ReadUint16(r); err != nil {
        return fmt.Errorf("failed to read message version: %s", err)
    } else if h.Version != Version {
        return fmt.Errorf("invalid message version: %s", h.Version)
    }
    if h.Type, err = basic.ReadUint8(r); err != nil {
        return fmt.Errorf("failed to read message type: %s", err)
    } else if h.Type == Unknown || h.Type > Cancelled {
        return fmt.Errorf("invalid message type: %s", h.Type)
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

type Message struct {
    Header Header
    Payload []byte
}

func (m Message) Write(w io.Writer) error {
    if err := m.Header.Write(w); err != nil {
        return fmt.Errorf("failed to write message header: %s", err)
    }
    bytes, err := w.Write(m.Payload)
    if (err != nil) {
        return fmt.Errorf("failed to write message payload: %s", err)
    } else if (bytes != int(m.Header.Size)) {
        return fmt.Errorf("failed to write message payload (%d instead of %d)", bytes, m.Header.Size)
    }
    return nil
}

func (m Message) Read(r io.Reader) error {
    if err := m.Header.Read(r); err != nil {
        return fmt.Errorf("failed to read message header: %s", err)
    }
    m.Payload = make([]byte, m.Header.Size)
    bytes, err := r.Read(m.Payload)
    if (err != nil) {
        return fmt.Errorf("failed to read message payload: %s", err)
    } else if (bytes != int(m.Header.Size)) {
        return fmt.Errorf("failed to read message payload (%d instead of %d)", bytes, m.Header.Size)
    }
    return nil
}

func NewMessage(header Header, payload []byte) Message {
    header.Size = uint32(len(payload))
    return Message { header, payload, }
}
