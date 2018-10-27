package basic_test

import (
	"bytes"
	"errors"
	. "github.com/lugu/qiloop/type/basic"
	"testing"
)

type EmptyReaderWriter struct{}

func (e EmptyReaderWriter) Read(p []byte) (n int, err error) {
	return 0, errors.New("nothing to read from")
}
func (e EmptyReaderWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("nothing to write to")
}

var readWriteError EmptyReaderWriter

func TestUint8(t *testing.T) {
	var values []uint8 = []uint8{
		0, 1, 2, 10, 30, 50, 127,
	}
	for _, value := range values {
		buf := bytes.NewBuffer(make([]byte, 0))
		if err := WriteUint8(value, buf); err != nil {
			t.Errorf("write error")
		}
		if result, err := ReadUint8(buf); err != nil {
			t.Errorf("read error")
		} else if result != value {
			t.Errorf("serialization error")
		}
		if err := WriteUint8(value, readWriteError); err == nil {
			t.Errorf("missing write error")
		}
		if _, err := ReadUint8(readWriteError); err == nil {
			t.Errorf("missingread error")
		}
	}
}

func TestUint16(t *testing.T) {
	var values []uint16 = []uint16{
		0, 1, 2, 10, 30, 50, 127,
		1<<11 + 127, 1 << 15,
	}
	for _, value := range values {
		buf := bytes.NewBuffer(make([]byte, 0))
		if err := WriteUint16(value, buf); err != nil {
			t.Errorf("write error")
		}
		if result, err := ReadUint16(buf); err != nil {
			t.Errorf("read error")
		} else if result != value {
			t.Errorf("serialization error")
		}
		if err := WriteUint16(value, readWriteError); err == nil {
			t.Errorf("missing write error")
		}
		if _, err := ReadUint16(readWriteError); err == nil {
			t.Errorf("missingread error")
		}
	}
}

func TestUint32(t *testing.T) {
	var values []uint32 = []uint32{
		0, 1, 2, 10, 30, 50, 127,
		1<<11 + 127, 1 << 16,
		1<<21 + 321, 1 << 28, 1 << 31,
	}
	for _, value := range values {
		buf := bytes.NewBuffer(make([]byte, 0))
		if err := WriteUint32(value, buf); err != nil {
			t.Errorf("write error")
		}
		if result, err := ReadUint32(buf); err != nil {
			t.Errorf("read error")
		} else if result != value {
			t.Errorf("serialization error")
		}
		if err := WriteUint32(value, readWriteError); err == nil {
			t.Errorf("missing write error")
		}
		if _, err := ReadUint32(readWriteError); err == nil {
			t.Errorf("missingread error")
		}
	}
}

func TestUint64(t *testing.T) {
	var values []uint64 = []uint64{
		0, 1, 2, 10, 30, 50, 127,
		1<<11 + 127, 1 << 16,
		1<<21 + 641, 1 << 28, 1 << 31,
		1<<51 + 641, 1 << 60, 1 << 63,
	}
	for _, value := range values {
		buf := bytes.NewBuffer(make([]byte, 0))
		if err := WriteUint64(value, buf); err != nil {
			t.Errorf("write error")
		}
		if result, err := ReadUint64(buf); err != nil {
			t.Errorf("read error")
		} else if result != value {
			t.Errorf("serialization error")
		}
		if err := WriteUint64(value, readWriteError); err == nil {
			t.Errorf("missing write error")
		}
		if _, err := ReadUint64(readWriteError); err == nil {
			t.Errorf("missingread error")
		}
	}
}

func TestInt32(t *testing.T) {
	var values []int32 = []int32{
		0, 1, 2, 10, 30, 50, 127,
		1<<11 + 127, 1 << 16,
		1<<21 + 321, 1 << 28, 1 << 30,
		-(1 << 30), -12, -1, -0, -(1 << 16),
	}
	for _, value := range values {
		buf := bytes.NewBuffer(make([]byte, 0))
		if err := WriteInt32(value, buf); err != nil {
			t.Errorf("write error")
		}
		if result, err := ReadInt32(buf); err != nil {
			t.Errorf("read error")
		} else if result != value {
			t.Errorf("serialization error")
		}
		if err := WriteInt32(value, readWriteError); err == nil {
			t.Errorf("missing write error")
		}
		if _, err := ReadInt32(readWriteError); err == nil {
			t.Errorf("missingread error")
		}
	}
}

func TestInt64(t *testing.T) {
	var values []int64 = []int64{
		0, 1, 2, 10, 30, 50, 127,
		1<<11 + 127, 1 << 16,
		1<<21 + 641, 1 << 28, 1 << 30,
		-(1 << 30), -12, -1, -0, -(1 << 16),
		-(1 << 40) + 12, -(1 << 63) + 12,
	}
	for _, value := range values {
		buf := bytes.NewBuffer(make([]byte, 0))
		if err := WriteInt64(value, buf); err != nil {
			t.Errorf("write error")
		}
		if result, err := ReadInt64(buf); err != nil {
			t.Errorf("read error")
		} else if result != value {
			t.Errorf("serialization error")
		}
		if err := WriteInt64(value, readWriteError); err == nil {
			t.Errorf("missing write error")
		}
		if _, err := ReadInt64(readWriteError); err == nil {
			t.Errorf("missingread error")
		}
	}
}

func TestFloat32(t *testing.T) {
	var values []float32 = []float32{
		0.0, 1.0, 2.1, 10.3, 1 / 3, -1.5, -1 / 6,
	}
	for _, value := range values {
		buf := bytes.NewBuffer(make([]byte, 0))
		if err := WriteFloat32(value, buf); err != nil {
			t.Errorf("write error")
		}
		if result, err := ReadFloat32(buf); err != nil {
			t.Errorf("read error")
		} else if result != value {
			t.Errorf("serialization error")
		}
		if err := WriteFloat32(value, readWriteError); err == nil {
			t.Errorf("missing write error")
		}
		if _, err := ReadFloat32(readWriteError); err == nil {
			t.Errorf("missingread error")
		}
	}
}

func TestFloat64(t *testing.T) {
	var values []float64 = []float64{
		0.0, 1.0, 2.1, 10.3, 1 / 3, -1.5, -1 / 6,
	}
	for _, value := range values {
		buf := bytes.NewBuffer(make([]byte, 0))
		if err := WriteFloat64(value, buf); err != nil {
			t.Errorf("write error")
		}
		if result, err := ReadFloat64(buf); err != nil {
			t.Errorf("read error")
		} else if result != value {
			t.Errorf("serialization error")
		}
		if err := WriteFloat64(value, readWriteError); err == nil {
			t.Errorf("missing write error")
		}
		if _, err := ReadFloat64(readWriteError); err == nil {
			t.Errorf("missingread error")
		}
	}
}

func TestBool(t *testing.T) {
	var values []bool = []bool{true, false}
	for _, value := range values {
		buf := bytes.NewBuffer(make([]byte, 0))
		if err := WriteBool(value, buf); err != nil {
			t.Errorf("write error")
		}
		if result, err := ReadBool(buf); err != nil {
			t.Errorf("read error")
		} else if result != value {
			t.Errorf("serialization error")
		}
		if err := WriteBool(value, readWriteError); err == nil {
			t.Errorf("missing write error")
		}
		if _, err := ReadBool(readWriteError); err == nil {
			t.Errorf("missingread error")
		}
	}
}

func TestString(t *testing.T) {
	var values []string = []string{
		"", `'"\{}\0 0`, "test", `Ã `,
		`Ã`, `some longer string  
		with returns
		`,
	}
	for _, value := range values {
		buf := bytes.NewBuffer(make([]byte, 0))
		if err := WriteString(value, buf); err != nil {
			t.Errorf("write error")
		}
		if result, err := ReadString(buf); err != nil {
			t.Errorf("read error")
		} else if result != value {
			t.Errorf("serialization error")
		}
		if err := WriteString(value, readWriteError); err == nil {
			t.Errorf("missing write error")
		}
		if _, err := ReadString(readWriteError); err == nil {
			t.Errorf("missingread error")
		}
	}
}
