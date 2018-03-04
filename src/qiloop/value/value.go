package value

import (
	"fmt"
	"io"
	"qiloop/basic"
)

type Value interface {
	Signature() string
	Write(w io.Writer) error
}

func NewValue(r io.Reader) (Value, error) {
	s, err := basic.ReadString(r)
	if err != nil {
		return nil, fmt.Errorf("value signature: %s", err)
	}
	switch s {
	case "i":
		return newInt(r)
	case "I":
		return newInt(r)
	case "L":
		return newLong(r)
	case "s":
		return newString(r)
	case "b":
		return newBool(r)
	case "f":
		return newFloat(r)
	default:
		return nil, fmt.Errorf("unsuported signature: %s", s)
	}
}

type BoolValue bool

func Bool(b bool) Value {
	return BoolValue(b)
}

func newBool(r io.Reader) (Value, error) {
	b, err := basic.ReadBool(r)
	return Bool(b), err
}

func (b BoolValue) Signature() string {
	return "b"
}

func (b BoolValue) Write(w io.Writer) error {
	if err := basic.WriteString(b.Signature(), w); err != nil {
		return err
	}
	return basic.WriteBool(b.Value(), w)
}

func (b BoolValue) Value() bool {
	return bool(b)
}

type IntValue uint32

func Int(i uint32) Value {
	return IntValue(i)
}

func newInt(r io.Reader) (Value, error) {
	i, err := basic.ReadUint32(r)
	return Int(i), err
}

func (i IntValue) Signature() string {
	return "I"
}

func (i IntValue) Write(w io.Writer) error {
	if err := basic.WriteString(i.Signature(), w); err != nil {
		return err
	}
	return basic.WriteUint32(i.Value(), w)
}

func (i IntValue) Value() uint32 {
	return uint32(i)
}

type LongValue uint64

func Long(l uint64) Value {
	return LongValue(l)
}

func newLong(r io.Reader) (Value, error) {
	l, err := basic.ReadUint64(r)
	return Long(l), err
}

func (l LongValue) Signature() string {
	return "L"
}

func (l LongValue) Write(w io.Writer) error {
	if err := basic.WriteString(l.Signature(), w); err != nil {
		return err
	}
	return basic.WriteUint64(l.Value(), w)
}

func (l LongValue) Value() uint64 {
	return uint64(l)
}

type FloatValue float32

func Float(f float32) Value {
	return FloatValue(f)
}

func newFloat(r io.Reader) (Value, error) {
	f, err := basic.ReadFloat32(r)
	return Float(f), err
}

func (f FloatValue) Signature() string {
	return "f"
}

func (f FloatValue) Write(w io.Writer) error {
	if err := basic.WriteString(f.Signature(), w); err != nil {
		return err
	}
	return basic.WriteFloat32(f.Value(), w)
}

func (f FloatValue) Value() float32 {
	return float32(f)
}

type StringValue string

func String(s string) Value {
	return StringValue(s)
}

func newString(r io.Reader) (Value, error) {
	s, err := basic.ReadString(r)
	return String(s), err
}

func (s StringValue) Signature() string {
	return "s"
}

func (s StringValue) Write(w io.Writer) error {
	if err := basic.WriteString(s.Signature(), w); err != nil {
		return err
	}
	return basic.WriteString(s.Value(), w)
}

func (s StringValue) Value() string {
	return string(s)
}
