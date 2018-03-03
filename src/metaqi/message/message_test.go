package message_test

import (
	"bytes"
	"metaqi/message"
	"os"
	"path/filepath"
	"testing"
)

func helpParseHeader(t *testing.T, filename string, expected message.Header) {
	var m message.Message
	path := filepath.Join("testdata", filename)
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = m.Read(file); err != nil {
		t.Error(err)
	}
	if m.Header != expected {
		t.Errorf("expected %#v, got %#v", expected, m.Header)
	}
}

func TestParseCallHeader(t *testing.T) {
	filename := "header-call-authenticate.bin"
	expected := message.Header{
		0x42dead42,
		3,   // id
		110, // size
		0,   // verison
		1,   // type call
		0,   // flags
		0,   // service
		0,   // object
		8,   // authenticate action
	}
	helpParseHeader(t, filename, expected)
}

func TestParseReplyHeader(t *testing.T) {
	filename := "header-reply-authenticate.bin"
	expected := message.Header{
		0x42dead42,
		3,   // id
		138, // size
		0,   // verison
		2,   // type call
		0,   // flags
		0,   // service
		0,   // object
		8,   // authenticate action
	}
	helpParseHeader(t, filename, expected)
}

func TestConstants(t *testing.T) {
	if 0x42dead42 != message.Magic {
		t.Error("invalid magic definition")
	}
	if 0 != message.Version {
		t.Error("invalid version definition")
	}
	if 1 != message.Call {
		t.Error("invalid call definition")
	}
	if 2 != message.Reply {
		t.Error("invalid call definition")
	}
	if 3 != message.Error {
		t.Error("invalid error definition")
	}
}

func TestMessageConstructor(t *testing.T) {
	h := message.NewHeader(message.Call, 1, 2, 3, 4)
	m := message.NewMessage(h, make([]byte, 99))

	if m.Header.Magic != message.Magic {
		t.Errorf("invalid magic: %d", m.Header.Magic)
	}
	if m.Header.Version != message.Version {
		t.Errorf("invalid version: %d", m.Header.Version)
	}
	if m.Header.Type != message.Call {
		t.Errorf("invalid type: %d", m.Header.Type)
	}
	if m.Header.Flags != 0 {
		t.Errorf("invalid flags: %d", m.Header.Flags)
	}
	if m.Header.Service != 1 {
		t.Errorf("invalid service: %d", m.Header.Service)
	}
	if m.Header.Object != 2 {
		t.Errorf("invalid service: %d", m.Header.Object)
	}
	if m.Header.Action != 3 {
		t.Errorf("invalid service: %d", m.Header.Action)
	}
	if m.Header.Id != 4 {
		t.Errorf("invalid id: %d", m.Header.Id)
	}
	if m.Header.Size != 99 {
		t.Errorf("invalid size: %d", m.Header.Size)
	}
}

func TestWriteReadMessage(t *testing.T) {
	h := message.NewHeader(message.Call, 1, 2, 3, 4)
	input := message.NewMessage(h, make([]byte, 99))
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := input.Write(buf); err != nil {
		t.Errorf("failed to write message: %s", err)
	}
	var output message.Message
	if err := output.Read(buf); err != nil {
		t.Errorf("failed to read message: %s", err)
	}
	if input.Header != output.Header {
		t.Errorf("expected %#v, got %#v", input.Header, output.Header)
	}
}
