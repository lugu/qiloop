package value

import (
	"fmt"
	"github.com/lugu/qiloop/type/basic"
	"io"
)

// Value represents a value whose type in unknown at compile time. The
// value can be an integer, a float, a boolean, a long or a string.
// When serialized, the signature of the true type is sent followed by
// the actual value.
type Value interface {
	signature() string
	Write(w io.Writer) error
}

// NewValue reads a value from a reader. The value is constructed in
// two times: first the signature of the value is read from the
// reader, then depending on the actual type, the value is read.
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

// BoolValue represents a Value of a boolean.
type BoolValue bool

// Bool constructs a Value.
func Bool(b bool) Value {
	return BoolValue(b)
}

func newBool(r io.Reader) (Value, error) {
	b, err := basic.ReadBool(r)
	return Bool(b), err
}

func (b BoolValue) signature() string {
	return "b"
}

func (b BoolValue) Write(w io.Writer) error {
	if err := basic.WriteString(b.signature(), w); err != nil {
		return err
	}
	return basic.WriteBool(b.Value(), w)
}

// Value returns the actual value.
func (b BoolValue) Value() bool {
	return bool(b)
}

// IntValue represents a Value of an uint32.
type IntValue uint32

// Int constructs a Value.
func Int(i uint32) Value {
	return IntValue(i)
}

func newInt(r io.Reader) (Value, error) {
	i, err := basic.ReadUint32(r)
	return Int(i), err
}

func (i IntValue) signature() string {
	return "I"
}

func (i IntValue) Write(w io.Writer) error {
	if err := basic.WriteString(i.signature(), w); err != nil {
		return err
	}
	return basic.WriteUint32(i.Value(), w)
}

// Value returns the actual value
func (i IntValue) Value() uint32 {
	return uint32(i)
}

// LongValue represents a Value of a uint64.
type LongValue uint64

// Long constructs a Value.
func Long(l uint64) Value {
	return LongValue(l)
}

func newLong(r io.Reader) (Value, error) {
	l, err := basic.ReadUint64(r)
	return Long(l), err
}

func (l LongValue) signature() string {
	return "L"
}

func (l LongValue) Write(w io.Writer) error {
	if err := basic.WriteString(l.signature(), w); err != nil {
		return err
	}
	return basic.WriteUint64(l.Value(), w)
}

// Value returns the actual value
func (l LongValue) Value() uint64 {
	return uint64(l)
}

// FloatValue represents a Value of a float32.
type FloatValue float32

// Float contructs a Value.
func Float(f float32) Value {
	return FloatValue(f)
}

func newFloat(r io.Reader) (Value, error) {
	f, err := basic.ReadFloat32(r)
	return Float(f), err
}

func (f FloatValue) signature() string {
	return "f"
}

func (f FloatValue) Write(w io.Writer) error {
	if err := basic.WriteString(f.signature(), w); err != nil {
		return err
	}
	return basic.WriteFloat32(f.Value(), w)
}

// Value returns the actual value
func (f FloatValue) Value() float32 {
	return float32(f)
}

// StringValue represents a Value of a string.
type StringValue string

// String constructs a Value.
func String(s string) Value {
	return StringValue(s)
}

func newString(r io.Reader) (Value, error) {
	s, err := basic.ReadString(r)
	return String(s), err
}

func (s StringValue) signature() string {
	return "s"
}

func (s StringValue) Write(w io.Writer) error {
	if err := basic.WriteString(s.signature(), w); err != nil {
		return err
	}
	return basic.WriteString(s.Value(), w)
}

// Value returns the actual value
func (s StringValue) Value() string {
	return string(s)
}
