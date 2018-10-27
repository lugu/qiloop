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
	case "c":
		return newInt8(r)
	case "C":
		return newUint8(r)
	case "w":
		return newInt16(r)
	case "W":
		return newUint16(r)
	case "i":
		return newInt(r)
	case "I":
		return newUint(r)
	case "l":
		return newLong(r)
	case "L":
		return newUlong(r)
	case "s":
		return newString(r)
	case "b":
		return newBool(r)
	case "f":
		return newFloat(r)
	case "[m]":
		return newList(r)
	case "r":
		return newRaw(r)
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

// Uint8Value represents a Value of an uint8.
type Uint8Value uint8

// Int constructs a Value. FIXME: Int shall be int32
func Uint8(i uint8) Value {
	return Uint8Value(i)
}

func newUint8(r io.Reader) (Value, error) {
	i, err := basic.ReadUint8(r)
	return Uint8(i), err
}

func (i Uint8Value) signature() string {
	return "I"
}

func (i Uint8Value) Write(w io.Writer) error {
	if err := basic.WriteString(i.signature(), w); err != nil {
		return err
	}
	return basic.WriteUint8(i.Value(), w)
}

// Value returns the actual value
func (i Uint8Value) Value() uint8 {
	return uint8(i)
}

// Int8Value represents a Value of an int8.
type Int8Value int8

// Int constructs a Value. FIXME: Int shall be int32
func Int8(i int8) Value {
	return Int8Value(i)
}

func newInt8(r io.Reader) (Value, error) {
	i, err := basic.ReadInt8(r)
	return Int8(i), err
}

func (i Int8Value) signature() string {
	return "i"
}

func (i Int8Value) Write(w io.Writer) error {
	if err := basic.WriteString(i.signature(), w); err != nil {
		return err
	}
	return basic.WriteInt8(i.Value(), w)
}

// Value returns the actual value
func (i Int8Value) Value() int8 {
	return int8(i)
}

// Uint16Value represents a Value of an uint16.
type Uint16Value uint16

// Int constructs a Value. FIXME: Int shall be int32
func Uint16(i uint16) Value {
	return Uint16Value(i)
}

func newUint16(r io.Reader) (Value, error) {
	i, err := basic.ReadUint16(r)
	return Uint16(i), err
}

func (i Uint16Value) signature() string {
	return "W"
}

func (i Uint16Value) Write(w io.Writer) error {
	if err := basic.WriteString(i.signature(), w); err != nil {
		return err
	}
	return basic.WriteUint16(i.Value(), w)
}

// Value returns the actual value
func (i Uint16Value) Value() uint16 {
	return uint16(i)
}

// Int16Value represents a Value of an int16.
type Int16Value int16

// Int constructs a Value. FIXME: Int shall be int32
func Int16(i int16) Value {
	return Int16Value(i)
}

func newInt16(r io.Reader) (Value, error) {
	i, err := basic.ReadInt16(r)
	return Int16(i), err
}

func (i Int16Value) signature() string {
	return "w"
}

func (i Int16Value) Write(w io.Writer) error {
	if err := basic.WriteString(i.signature(), w); err != nil {
		return err
	}
	return basic.WriteInt16(i.Value(), w)
}

// Value returns the actual value
func (i Int16Value) Value() int16 {
	return int16(i)
}

// UintValue represents a Value of an uint32.
type UintValue uint32

// Int constructs a Value. FIXME: Int shall be int32
func Uint(i uint32) Value {
	return UintValue(i)
}

func newUint(r io.Reader) (Value, error) {
	i, err := basic.ReadUint32(r)
	return Uint(i), err
}

func (i UintValue) signature() string {
	return "I"
}

func (i UintValue) Write(w io.Writer) error {
	if err := basic.WriteString(i.signature(), w); err != nil {
		return err
	}
	return basic.WriteUint32(i.Value(), w)
}

// Value returns the actual value
func (i UintValue) Value() uint32 {
	return uint32(i)
}

// IntValue represents a Value of an uint32.
type IntValue int32

// Int constructs a Value. FIXME: Int shall be int32
func Int(i int32) Value {
	return IntValue(i)
}

func newInt(r io.Reader) (Value, error) {
	i, err := basic.ReadInt32(r)
	return Int(i), err
}

func (i IntValue) signature() string {
	return "i"
}

func (i IntValue) Write(w io.Writer) error {
	if err := basic.WriteString(i.signature(), w); err != nil {
		return err
	}
	return basic.WriteInt32(i.Value(), w)
}

// Value returns the actual value
func (i IntValue) Value() int32 {
	return int32(i)
}

// UlongValue represents a Value of a uint64.
type UlongValue uint64

// Ulong constructs a Value.
func Ulong(l uint64) Value {
	return UlongValue(l)
}

func newUlong(r io.Reader) (Value, error) {
	l, err := basic.ReadUint64(r)
	return Ulong(l), err
}

func (l UlongValue) signature() string {
	return "L"
}

func (l UlongValue) Write(w io.Writer) error {
	if err := basic.WriteString(l.signature(), w); err != nil {
		return err
	}
	return basic.WriteUint64(l.Value(), w)
}

// Value returns the actual value
func (l UlongValue) Value() uint64 {
	return uint64(l)
}

// LongValue represents a Value of a uint64.
type LongValue int64

// Long constructs a Value.
func Long(l int64) Value {
	return LongValue(l)
}

func newLong(r io.Reader) (Value, error) {
	l, err := basic.ReadInt64(r)
	return Long(l), err
}

func (l LongValue) signature() string {
	return "l"
}

func (l LongValue) Write(w io.Writer) error {
	if err := basic.WriteString(l.signature(), w); err != nil {
		return err
	}
	return basic.WriteInt64(l.Value(), w)
}

// Value returns the actual value
func (l LongValue) Value() int64 {
	return int64(l)
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

type ListValue []Value

func newList(r io.Reader) (Value, error) {
	size, err := basic.ReadUint32(r)
	if err != nil {
		return nil, err
	}
	list := make([]Value, size)
	for i := range list {
		list[i], err = NewValue(r)
		if err != nil {
			return nil, err
		}
	}
	return ListValue(list), err
}

// List constructs a Value.
func List(l []Value) Value {
	return ListValue(l)
}

func (l ListValue) signature() string {
	return "[m]"
}

func (l ListValue) Write(w io.Writer) error {
	if err := basic.WriteString(l.signature(), w); err != nil {
		return err
	}
	if err := basic.WriteUint32(uint32(len(l)), w); err != nil {
		return err
	}
	for _, v := range l {
		if err := v.Write(w); err != nil {
			return err
		}
	}
	return nil
}

// Value returns the actual value
func (l ListValue) Value() []Value {
	return []Value(l)
}

type RawValue []byte

const (
	RawValueMaxSize = 10 * 1024 * 1024
)

func newRaw(r io.Reader) (Value, error) {
	size, err := basic.ReadUint32(r)
	if err != nil {
		return nil, err
	}
	if size > RawValueMaxSize {
		return nil, fmt.Errorf("raw value too large (%d bytes)", size)
	}
	buf := make([]byte, size)
	count := 0
	for count < int(size) {
		read, err := r.Read(buf[count:])
		count += read
		if err != nil && count < int(size) {
			return nil, fmt.Errorf("raw read: %d instead of %d: %s",
				count, size, err)
		}
	}
	return RawValue(buf), nil
}

// Raw constructs a Value.
func Raw(b []byte) Value {
	return RawValue(b)
}

func (b RawValue) signature() string {
	return "r"
}

func (b RawValue) Write(w io.Writer) error {
	if err := basic.WriteString(b.signature(), w); err != nil {
		return err
	}
	if err := basic.WriteUint32(uint32(len(b)), w); err != nil {
		return err
	}
	count := 0
	for count < int(len(b)) {
		written, err := w.Write(b[count:])
		count += written
		if err != nil && count < int(len(b)) {
			return fmt.Errorf("raw write: %d instead of %d: %s",
				count, len(b), err)
		}
	}
	return nil
}

// Value returns the actual value
func (b RawValue) Value() []byte {
	return []byte(b)
}
