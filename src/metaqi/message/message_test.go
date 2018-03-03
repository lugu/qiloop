package message_test

import (
    "bytes"
    "path/filepath"
    "os"
    "metaqi/message"
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

func testParseCallHeader(t *testing.T) {
    filename := "header-call-authenticate.bin"
    expected := message.Header {
        0x42dead42,
        3, // id
        110, // size
        0, // verison
        1, // type call
        0, // flags
        0, // service
        0, // object
        8, // authenticate action
    }
    helpParseHeader(t, filename, expected)
}

func testParseReplyHeader(t *testing.T) {
    filename := "header-reply-authenticate.bin"
    expected := message.Header {
        0x42dead42,
        3, // id
        138, // size
        0, // verison
        2, // type call
        0, // flags
        0, // service
        0, // object
        8, // authenticate action
    }
    helpParseHeader(t, filename, expected)
}

func testConstants(t *testing.T) {
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

func testMessageConstructor(t *testing.T) {
    h := message.NewHeader(message.Call, 1, 2, 3, 4)
    m := message.NewMessage(h, make([]byte, 99))

    if m.Header.Magic != message.Magic {
        t.Error("invalid magic")
    }
    if m.Header.Version != message.Version {
        t.Error("invalid version")
    }
    if m.Header.Id == 4 {
        t.Error("invalid id")
    }
    if m.Header.Type != message.Call {
        t.Error("invalid type")
    }
    if m.Header.Flags != 0 {
        t.Error("invalid flags")
    }
    if m.Header.Service != 1 {
        t.Error("invalid service")
    }
    if m.Header.Object != 2 {
        t.Error("invalid service")
    }
    if m.Header.Action != 3 {
        t.Error("invalid service")
    }
    if m.Header.Size != 99 {
        t.Error("invalid size")
    }
}

func testWriteReadMessage(t *testing.T) {
    h := message.NewHeader(message.Call, 1, 2, 3, 4)
    input := message.NewMessage(h, make([]byte, 99))
    buf := bytes.NewBuffer(make([]byte, 1))
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
