package value

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/basic"
)

const (
	rawValueMaxSize  = 10 * 1024 * 1024
	listValueMaxSize = 4096

	ObjectReferenceSignature = "(({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>II)<ObjectReference,metaObject,serviceID,objectID>"
)

var (
	// ErrListValueTooLong is returned when reading a list value.
	ErrListValueTooLong = errors.New("list value too long")
	// ErrRawValueTooLong is returned when reading a raw value.
	ErrRawValueTooLong = errors.New("raw value too long")
)

// Value represents a value whose type in unknown at compile time. The
// value can be an integer, a float, a boolean, a long or a string.
// When serialized, the signature of the true type is sent followed by
// the actual value.
type Value interface {
	Signature() string
	Write(w io.Writer) error
}

// Bytes returns the content of the value (i.e. without the
// signature).
func Bytes(v Value) []byte {
	var buf bytes.Buffer
	v.Write(&buf)
	basic.ReadString(&buf)
	return buf.Bytes()
}

// NewValue reads a value from a reader. The value is constructed in
// two times: first the signature of the value is read from the
// reader, then depending on the actual type, the value is read.
func NewValue(r io.Reader) (Value, error) {
	s, err := basic.ReadString(r)
	if err != nil {
		return nil, fmt.Errorf("value signature: %s", err)
	}
	solve := map[string]func(io.Reader) (Value, error){
		"c":   newInt8,
		"C":   newUint8,
		"w":   newInt16,
		"W":   newUint16,
		"i":   newInt,
		"I":   newUint,
		"l":   newLong,
		"L":   newUlong,
		"s":   newString,
		"b":   newBool,
		"f":   newFloat,
		"[m]": newList,
		"r":   newRaw,
		"v":   newVoid,
		"m":   NewValue,
	}
	f, ok := solve[s]
	if ok {
		return f(r)
	}
	return newOpaque(s, r)
}

// OpaqueValue represents a value using a signature and the data.
type OpaqueValue struct {
	sig  string
	data []byte
}

// Opaque creates a Value (an OpaqueValue).
func Opaque(signature string, data []byte) Value {
	return &OpaqueValue{
		sig:  signature,
		data: data,
	}
}

// Signature returns the signature of the value.
func (o *OpaqueValue) Signature() string {
	return o.sig
}

// Write the value.
func (o *OpaqueValue) Write(w io.Writer) error {
	if err := basic.WriteString(o.sig, w); err != nil {
		return err
	}
	err := basic.WriteN(w, o.data, len(o.data))
	if err != nil {
		return fmt.Errorf("write value: %s", err)
	}
	return nil
}

func newOpaque(sig string, r io.Reader) (Value, error) {
	if sig == "o" {
		return newOpaque(ObjectReferenceSignature, r)
	}
	reader, err := signature.MakeReader(sig)
	if err != nil {
		return nil, fmt.Errorf("Invalid signature %s: %s", sig, err)
	}
	data, err := reader.Read(r)
	if err != nil {
		return nil, fmt.Errorf("Failed to read value %s: %s", sig, err)
	}
	return &OpaqueValue{
		sig:  sig,
		data: data,
	}, nil
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

// Signature returns the boolean signature "b".
func (b BoolValue) Signature() string {
	return "b"
}

func (b BoolValue) Write(w io.Writer) error {
	if err := basic.WriteString(b.Signature(), w); err != nil {
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

// Uint8 constructs a Value corresponding to an unsigned char.
func Uint8(i uint8) Value {
	return Uint8Value(i)
}

func newUint8(r io.Reader) (Value, error) {
	i, err := basic.ReadUint8(r)
	return Uint8(i), err
}

// Signature returns the signature of an unsigned char.
func (i Uint8Value) Signature() string {
	return "C"
}

func (i Uint8Value) Write(w io.Writer) error {
	if err := basic.WriteString(i.Signature(), w); err != nil {
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

// Int8 constructs a Value.
func Int8(i int8) Value {
	return Int8Value(i)
}

func newInt8(r io.Reader) (Value, error) {
	i, err := basic.ReadInt8(r)
	return Int8(i), err
}

// Signature returns the signature of an signed char ("c").
func (i Int8Value) Signature() string {
	return "c"
}

func (i Int8Value) Write(w io.Writer) error {
	if err := basic.WriteString(i.Signature(), w); err != nil {
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

// Uint16 constructs a Value.
func Uint16(i uint16) Value {
	return Uint16Value(i)
}

func newUint16(r io.Reader) (Value, error) {
	i, err := basic.ReadUint16(r)
	return Uint16(i), err
}

// Signature returns "W"
func (i Uint16Value) Signature() string {
	return "W"
}

func (i Uint16Value) Write(w io.Writer) error {
	if err := basic.WriteString(i.Signature(), w); err != nil {
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

// Int16 constructs a Value.
func Int16(i int16) Value {
	return Int16Value(i)
}

func newInt16(r io.Reader) (Value, error) {
	i, err := basic.ReadInt16(r)
	return Int16(i), err
}

// Signature returns "w".
func (i Int16Value) Signature() string {
	return "w"
}

func (i Int16Value) Write(w io.Writer) error {
	if err := basic.WriteString(i.Signature(), w); err != nil {
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

// Uint constructs a Value.
func Uint(i uint32) Value {
	return UintValue(i)
}

func newUint(r io.Reader) (Value, error) {
	i, err := basic.ReadUint32(r)
	return Uint(i), err
}

// Signature returns "I".
func (i UintValue) Signature() string {
	return "I"
}

func (i UintValue) Write(w io.Writer) error {
	if err := basic.WriteString(i.Signature(), w); err != nil {
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

// Int constructs a Value.
func Int(i int32) Value {
	return IntValue(i)
}

func newInt(r io.Reader) (Value, error) {
	i, err := basic.ReadInt32(r)
	return Int(i), err
}

// Signature returns "i".
func (i IntValue) Signature() string {
	return "i"
}

func (i IntValue) Write(w io.Writer) error {
	if err := basic.WriteString(i.Signature(), w); err != nil {
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

// Signature returns "L".
func (l UlongValue) Signature() string {
	return "L"
}

func (l UlongValue) Write(w io.Writer) error {
	if err := basic.WriteString(l.Signature(), w); err != nil {
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

// Signature returns "l"
func (l LongValue) Signature() string {
	return "l"
}

func (l LongValue) Write(w io.Writer) error {
	if err := basic.WriteString(l.Signature(), w); err != nil {
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

// Signature returns "f".
func (f FloatValue) Signature() string {
	return "f"
}

func (f FloatValue) Write(w io.Writer) error {
	if err := basic.WriteString(f.Signature(), w); err != nil {
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

// Signature returns "s".
func (s StringValue) Signature() string {
	return "s"
}

func (s StringValue) Write(w io.Writer) error {
	if err := basic.WriteString(s.Signature(), w); err != nil {
		return err
	}
	return basic.WriteString(s.Value(), w)
}

// Value returns the actual value
func (s StringValue) Value() string {
	return string(s)
}

// ListValue represents a list of value.
type ListValue []Value

func newList(r io.Reader) (Value, error) {
	size, err := basic.ReadUint32(r)
	if err != nil {
		return nil, err
	}
	if size > listValueMaxSize {
		return nil, ErrListValueTooLong
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

// Signature returns "[m]".
func (l ListValue) Signature() string {
	return "[m]"
}

func (l ListValue) Write(w io.Writer) error {
	if err := basic.WriteString(l.Signature(), w); err != nil {
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

// RawValue represents an array of byte.
type RawValue []byte

func newRaw(r io.Reader) (Value, error) {
	size, err := basic.ReadUint32(r)
	if err != nil {
		return nil, err
	}
	if size > rawValueMaxSize {
		return nil, ErrRawValueTooLong
	}
	buf := make([]byte, size)
	err = basic.ReadN(r, buf, int(size))
	if err != nil {
		return nil, fmt.Errorf("raw read: %s", err)
	}
	return RawValue(buf), nil
}

// Raw constructs a Value.
func Raw(b []byte) Value {
	return RawValue(b)
}

// Signature returns "r".
func (b RawValue) Signature() string {
	return "r"
}

func (b RawValue) Write(w io.Writer) error {
	if err := basic.WriteString(b.Signature(), w); err != nil {
		return err
	}
	if err := basic.WriteUint32(uint32(len(b)), w); err != nil {
		return err
	}
	err := basic.WriteN(w, b, len(b))
	if err != nil {
		return fmt.Errorf("raw write: %s", err)
	}
	return nil
}

// Value returns the actual value
func (b RawValue) Value() []byte {
	return []byte(b)
}

// VoidValue represents nothing.
type VoidValue struct{}

// Void returns an empty value
func Void() Value {
	return VoidValue{}
}

func newVoid(r io.Reader) (Value, error) {
	return VoidValue{}, nil
}

// Signature returns "v".
func (b VoidValue) Signature() string {
	return "v"
}

func (b VoidValue) Write(w io.Writer) error {
	if err := basic.WriteString(b.Signature(), w); err != nil {
		return err
	}
	return nil
}
