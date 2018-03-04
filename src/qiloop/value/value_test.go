package value_test

import (
	"bytes"
	"qiloop/value"
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

func helpParseValue(t *testing.T, b []byte, expected value.Value) {
	buf := bytes.NewBuffer(b)
	v, err := value.NewValue(buf)
	if err != nil {
		t.Errorf("failed to parse value: %s", err)
	}
	if !reflect.DeepEqual(v, expected) {
		t.Errorf("expected %#v, got %#v", expected, v)
	}
}

func TestParseBool(t *testing.T) {
	bytes := []byte{1, 0, 0, 0, 0x62, 1}
	helpParseValue(t, bytes, value.Bool(true))
}

func TestParseInt(t *testing.T) {
	bytes := []byte{1, 0, 0, 0, 0x49, 0xff, 0, 0, 0xee}
	helpParseValue(t, bytes, value.Int(0xee0000ff))
}

func TestParseLong(t *testing.T) {
	bytes := []byte{1, 0, 0, 0, 0x4c, 0xff, 0, 0, 0, 0, 0, 0, 0xee}
	helpParseValue(t, bytes, value.Long(0xee000000000000ff))
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
