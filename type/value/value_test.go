package value_test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/lugu/qiloop/type/value"
)

func TestValues(t *testing.T) {
	if value.Bool(true).(value.BoolValue).Value() != true {
		t.Errorf("value constructor error")
	}
	if value.Bool(false).(value.BoolValue).Value() != false {
		t.Errorf("value constructor error")
	}
	if value.Int(0).(value.IntValue).Value() != 0 {
		t.Errorf("value constructor error")
	}
	if value.Int(12).(value.IntValue).Value() != 12 {
		t.Errorf("value constructor error")
	}
	if value.Long(12).(value.LongValue).Value() != 12 {
		t.Errorf("value constructor error")
	}
	if value.String("aa").(value.StringValue).Value() != "aa" {
		t.Errorf("value constructor error")
	}
	if value.Float(1.33).(value.FloatValue).Value() != 1.33 {
		t.Errorf("value constructor error")
	}
}

type LimitedWriter struct {
	size int
}

func (b *LimitedWriter) Write(buf []byte) (int, error) {
	if len(buf) <= b.size {
		b.size -= len(buf)
		return len(buf), nil
	}
	oldSize := b.size
	b.size = 0
	return oldSize, io.EOF
}

func NewLimitedWriter(size int) io.Writer {
	return &LimitedWriter{
		size: size,
	}
}

func helpValueWrite(t *testing.T, expected value.Value) {
	var buf bytes.Buffer
	err := expected.Write(&buf)
	if err != nil {
		t.Fatalf("write: %s", err)
	}
	size := len(buf.Bytes())
	val, err := value.NewValue(&buf)
	if err != nil {
		t.Fatalf("read: %s\n%#v", err, expected)
	} else if !reflect.DeepEqual(val, expected) {
		t.Fatalf("error: val different from expected:\n%#v\n%#v",
			val, expected)
	}

	err = expected.Write(NewLimitedWriter(4))
	if err == nil {
		t.Errorf("shall not work")
	}
	err = expected.Write(NewLimitedWriter(size - 1))
	if err == nil {
		t.Errorf("shall not work")
	}
}

func TestValueWriteRead(t *testing.T) {
	helpValueWrite(t, value.Bool(true))
	helpValueWrite(t, value.Bool(false))
	helpValueWrite(t, value.Int8(0))
	helpValueWrite(t, value.Int8(-42))
	helpValueWrite(t, value.Uint8(0))
	helpValueWrite(t, value.Uint8(42))
	helpValueWrite(t, value.Int16(0))
	helpValueWrite(t, value.Int16(-42))
	helpValueWrite(t, value.Uint16(0))
	helpValueWrite(t, value.Uint16(42))
	helpValueWrite(t, value.Int(0))
	helpValueWrite(t, value.Int(42))
	helpValueWrite(t, value.Long(0))
	helpValueWrite(t, value.Long(42<<42))
	helpValueWrite(t, value.Float(-1.234))
	helpValueWrite(t, value.Float(0))
	helpValueWrite(t, value.String(""))
	helpValueWrite(t, value.String("keep testing"))
	helpValueWrite(t, value.Void())
	helpValueWrite(t, value.Raw([]byte{1, 2, 3}))
	helpValueWrite(t, value.List([]value.Value{
		value.String("abc"),
		value.Int(0),
		value.List([]value.Value{
			value.String("abc"),
			value.Int(0),
		}),
	}))
	helpValueWrite(t, value.Opaque("(b)", []byte{1}))
	helpValueWrite(t, value.Opaque("(b)<Foo,a>", []byte{1}))
	helpValueWrite(t, value.Opaque("(bi)<Foo,a,b>", []byte{1, 1, 2, 3, 4}))
	helpValueWrite(t, value.Opaque("(b(i)<Bar,b>)<Foo,a,b>", []byte{1, 1, 2, 3, 4}))
	helpValueWrite(t, value.Opaque("(s)<Foo,a>", []byte{3, 0, 0, 0, 'a', 'b', 'c'}))
	helpValueWrite(t, value.Opaque("[(s)<Foo,a>]", []byte{
		3, 0, 0, 0,
		3, 0, 0, 0, 'a', 'b', 'c',
		3, 0, 0, 0, 'a', 'b', 'c',
		3, 0, 0, 0, 'a', 'b', 'c',
	}))
	helpValueWrite(t, value.Opaque("[b]", []byte{
		0, 0, 0, 0,
	}))
	helpValueWrite(t, value.Opaque("[b]", []byte{
		5, 0, 0, 0,
		1, 1, 1, 1, 1,
	}))
	helpValueWrite(t, value.Opaque("{ib}", []byte{
		0, 0, 0, 0,
	}))
	helpValueWrite(t, value.Opaque("{ib}", []byte{
		1, 0, 0, 0,
		1, 0, 0, 1, 1,
	}))
	helpValueWrite(t, value.Opaque("{ib}", []byte{
		2, 0, 0, 0,
		1, 0, 0, 1, 0,
		2, 0, 0, 2, 1,
	}))
	helpValueWrite(t, value.Opaque("{ib}", []byte{
		6, 0, 0, 0,
		1, 0, 0, 1, 0,
		2, 0, 0, 2, 1,
		3, 0, 0, 2, 0,
		4, 0, 0, 2, 1,
		5, 0, 0, 2, 0,
		6, 0, 0, 2, 1,
	}))
	helpValueWrite(t, value.Opaque("{i(s)<Foo,a>}", []byte{
		4, 0, 0, 0,
		1, 0, 0, 0,
		3, 0, 0, 0, 'a', 'b', 'c',
		2, 0, 0, 0,
		3, 0, 0, 0, 'a', 'b', 'c',
		3, 0, 0, 0,
		3, 0, 0, 0, 'a', 'b', 'c',
		4, 0, 0, 0,
		3, 0, 0, 0, 'a', 'b', 'c',
	}))
}

func helpParseValue(t *testing.T, b []byte, expected value.Value) {
	buf := bytes.NewBuffer(b)
	v, err := value.NewValue(buf)
	if err != nil {
		t.Errorf("parse value: %s", err)
	}
	if !reflect.DeepEqual(v, expected) {
		t.Errorf("expected %#v, got %#v", expected, v)
	}
	var out bytes.Buffer
	v.Write(&out)
	actual := out.Bytes()
	if !reflect.DeepEqual(b, actual) {
		t.Errorf("write: expected %#v, got %#v", b, actual)
	}
}

func TestParseBool(t *testing.T) {
	bytes := []byte{1, 0, 0, 0, 0x62, 1}
	helpParseValue(t, bytes, value.Bool(true))
}

func TestParseUint(t *testing.T) {
	bytes := []byte{1, 0, 0, 0, 0x49, 0xff, 0, 0, 0xee}
	helpParseValue(t, bytes, value.Uint(0xee0000ff))
}

func TestParseInt(t *testing.T) {
	bytes := []byte{1, 0, 0, 0, 0x69, 0xff, 0, 0, 0x7f}
	helpParseValue(t, bytes, value.Int(0x7f0000ff))
}

func TestParseUlong(t *testing.T) {
	bytes := []byte{1, 0, 0, 0, 0x4c, 0xff, 0, 0, 0, 0, 0, 0, 0xee}
	helpParseValue(t, bytes, value.Ulong(0xee000000000000ff))
}

func TestParseLong(t *testing.T) {
	bytes := []byte{1, 0, 0, 0, 0x6c, 0xff, 0, 0, 0, 0, 0, 0, 0x77}
	helpParseValue(t, bytes, value.Long(0x77000000000000ff))
}

func TestParseString(t *testing.T) {
	bytes := []byte{1, 0, 0, 0, 0x73, 03, 0, 0, 0, 0x4c, 0x49, 0x62}
	helpParseValue(t, bytes, value.String("LIb"))
}

func TestParseFloat(t *testing.T) {
	// https://www.h-schmidt.net/FloatConverter/IEEE754.html
	bytes := []byte{1, 0, 0, 0, 0x66, 0, 0, 0xc0, 0x3f}
	helpParseValue(t, bytes, value.Float(1.5))
	bytes = []byte{1, 0, 0, 0, 0x66, 0, 0x3e, 0x1c, 0xc6}
	helpParseValue(t, bytes, value.Float(-9999.5))
}

func TestParseRawData(t *testing.T) {
	bytes := []byte{1, 0, 0, 0, 0x72, 3, 0, 0, 0, 0x61, 0x62, 0x63}
	helpParseValue(t, bytes, value.Raw([]byte{'a', 'b', 'c'}))
}

func TestParseListValue(t *testing.T) {
	bytes := []byte{
		3, 0, 0, 0, 0x5b, 0x6d, 0x5d, 2, 0, 0, 0, 1, 0, 0, 0,
		0x73, 03, 0, 0, 0, 0x4c, 0x49, 0x62, 1, 0, 0, 0, 0x62, 1,
	}
	helpParseValue(t, bytes, value.List([]value.Value{
		value.String("LIb"),
		value.Bool(true),
	}))
}

func TestParseVoidValue(t *testing.T) {
	bytes := []byte{1, 0, 0, 0, 0x76}
	helpParseValue(t, bytes, value.Void())
}

func TestBytes(t *testing.T) {

	if !reflect.DeepEqual(
		value.Bytes(value.Bool(true)),
		[]byte{1},
	) {
		t.Errorf("true is not one")
	}
	if !reflect.DeepEqual(
		value.Bytes(value.Bool(false)),
		[]byte{0},
	) {
		t.Errorf("false is not zero")
	}
	if !reflect.DeepEqual(
		value.Bytes(value.Int(8)),
		[]byte{8, 0, 0, 0},
	) {
		t.Errorf("8 is not 8")
	}
	if !reflect.DeepEqual(
		value.Bytes(value.Opaque("boo", []byte{1, 2, 3})),
		[]byte{1, 2, 3},
	) {
		t.Errorf("123 is not 123")
	}
	data := []byte{
		4, 0, 0, 0,
		1, 0, 0, 0,
		3, 0, 0, 0, 'a', 'b', 'c',
		2, 0, 0, 0,
		3, 0, 0, 0, 'a', 'b', 'c',
		3, 0, 0, 0,
		3, 0, 0, 0, 'a', 'b', 'c',
		4, 0, 0, 0,
		3, 0, 0, 0, 'a', 'b', 'c',
	}
	if !reflect.DeepEqual(
		value.Bytes(value.Opaque("{i(s)<Foo,a>}", data)),
		data,
	) {
		t.Errorf("data does not match")
	}
}
