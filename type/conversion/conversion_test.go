package conversion_test

import (
	"testing"

	"github.com/lugu/qiloop/type/conversion"
)

func newBool(b bool) *bool       { return &b }
func newString(s string) *string { return &s }

func newUint(n uint) *uint       { return &n }
func newUint8(n uint8) *uint8 { return &n }
func newUint16(n uint16) *uint16 { return &n }
func newUint32(n uint32) *uint32 { return &n }
func newUint64(n uint64) *uint64 { return &n }

func newInt(n int) *int       { return &n }
func newInt8(n int8) *int8 { return &n }
func newInt16(n int16) *int16 { return &n }
func newInt32(n int32) *int32 { return &n }
func newInt64(n int64) *int64 { return &n }

func helpTestConvert2(t *testing.T, a interface{}, b interface{}, v int64) {
	err := conversion.ConvertFrom(a, b)
	if err != nil {
		t.Error(err)
		return
	}
	i, ok := conversion.AsInt64(a)
	if !ok {
		panic("not an integer")
	}
	if i != v {
		t.Errorf("invalid conversion: %d sv %d", i, v)
	}
}

func helpTestConvert(t *testing.T, a interface{}) {

	values := []int64 {
		0, 1, 11, 111, 127,
	}
	for _, v := range values {
		helpTestConvert2(t, a, interface{}(newInt(int(v))), v)
		helpTestConvert2(t, a, int(v), v)
		helpTestConvert2(t, a, interface{}(newUint(uint(v))), v)
		helpTestConvert2(t, a, uint(v), v)
		helpTestConvert2(t, a, interface{}(newInt8(int8(v))), v)
		helpTestConvert2(t, a, int8(v), v)
		helpTestConvert2(t, a, interface{}(newUint8(uint8(v))), v)
		helpTestConvert2(t, a, uint8(v), v)
		helpTestConvert2(t, a, interface{}(newInt16(int16(v))), v)
		helpTestConvert2(t, a, int16(v), v)
		helpTestConvert2(t, a, interface{}(newUint16(uint16(v))), v)
		helpTestConvert2(t, a, uint16(v), v)
		helpTestConvert2(t, a, interface{}(newInt32(int32(v))), v)
		helpTestConvert2(t, a, int32(v), v)
		helpTestConvert2(t, a, interface{}(newUint32(uint32(v))), v)
		helpTestConvert2(t, a, uint32(v), v)
		helpTestConvert2(t, a, interface{}(newInt64(int64(v))), v)
		helpTestConvert2(t, a, int64(v), v)
		helpTestConvert2(t, a, interface{}(newUint64(uint64(v))), v)
		helpTestConvert2(t, a, uint64(v), v)
	}
}

func TestBasicConversion(t *testing.T) {
	helpTestConvert(t, interface{}(newInt(1)))
	helpTestConvert(t, interface{}(newInt16(1)))
	helpTestConvert(t, interface{}(newInt32(1)))
	helpTestConvert(t, interface{}(newInt64(1)))
}
