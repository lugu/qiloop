package value_test

import (
	"bytes"
	"github.com/lugu/qiloop/type/value"
	"reflect"
	"testing"
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

func helpValueWrite(t *testing.T, expected value.Value) {
	var buf bytes.Buffer
	err := expected.Write(&buf)
	if err != nil {
		t.Errorf("failed to write: %s", err)
	}
	val, err := value.NewValue(&buf)
	if err != nil {
		t.Errorf("failed to read: %s", err)
	} else if !reflect.DeepEqual(val, expected) {
		t.Errorf("value constructor error")
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
}

func helpParseValue(t *testing.T, b []byte, expected value.Value) {
	buf := bytes.NewBuffer(b)
	v, err := value.NewValue(buf)
	if err != nil {
		t.Errorf("failed to parse value: %s", err)
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
