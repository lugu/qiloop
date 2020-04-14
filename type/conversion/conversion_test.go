package conversion_test

import (
	"reflect"
	"testing"

	"github.com/lugu/qiloop/type/conversion"
)

func newBool(b bool) *bool       { return &b }
func newString(s string) *string { return &s }

func newFloat32(f float32) *float32 { return &f }
func newFloat64(f float64) *float64 { return &f }

func newUint(n uint) *uint       { return &n }
func newUint8(n uint8) *uint8    { return &n }
func newUint16(n uint16) *uint16 { return &n }
func newUint32(n uint32) *uint32 { return &n }
func newUint64(n uint64) *uint64 { return &n }

func newInt(n int) *int       { return &n }
func newInt8(n int8) *int8    { return &n }
func newInt16(n int16) *int16 { return &n }
func newInt32(n int32) *int32 { return &n }
func newInt64(n int64) *int64 { return &n }

func helpTestConvertInt2(t *testing.T, a interface{}, b interface{}, i int64) {
	err := conversion.ConvertFrom(a, b)
	if err != nil {
		t.Error(err)
		return
	}
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	in, ok := conversion.AsInt64(v)
	if !ok {
		panic("not an integer")
	}
	if in != i {
		t.Errorf("invalid conversion: %d sv %d", in, i)
	}
}

func helpTestConvertInt(t *testing.T, a interface{}) {

	values := []int64{
		0, 1, 11, 111, 127,
	}
	for _, v := range values {
		helpTestConvertInt2(t, a, interface{}(newInt(int(v))), v)
		helpTestConvertInt2(t, a, int(v), v)
		helpTestConvertInt2(t, a, interface{}(newUint(uint(v))), v)
		helpTestConvertInt2(t, a, uint(v), v)
		helpTestConvertInt2(t, a, interface{}(newInt8(int8(v))), v)
		helpTestConvertInt2(t, a, int8(v), v)
		helpTestConvertInt2(t, a, interface{}(newUint8(uint8(v))), v)
		helpTestConvertInt2(t, a, uint8(v), v)
		helpTestConvertInt2(t, a, interface{}(newInt16(int16(v))), v)
		helpTestConvertInt2(t, a, int16(v), v)
		helpTestConvertInt2(t, a, interface{}(newUint16(uint16(v))), v)
		helpTestConvertInt2(t, a, uint16(v), v)
		helpTestConvertInt2(t, a, interface{}(newInt32(int32(v))), v)
		helpTestConvertInt2(t, a, int32(v), v)
		helpTestConvertInt2(t, a, interface{}(newUint32(uint32(v))), v)
		helpTestConvertInt2(t, a, uint32(v), v)
		helpTestConvertInt2(t, a, interface{}(newInt64(int64(v))), v)
		helpTestConvertInt2(t, a, int64(v), v)
		helpTestConvertInt2(t, a, interface{}(newUint64(uint64(v))), v)
		helpTestConvertInt2(t, a, uint64(v), v)
	}
}

func TestInt(t *testing.T) {
	helpTestConvertInt(t, interface{}(newInt(1)))
	helpTestConvertInt(t, interface{}(newInt16(1)))
	helpTestConvertInt(t, interface{}(newInt32(1)))
	helpTestConvertInt(t, interface{}(newInt64(1)))
}

func TestString(t *testing.T) {
	str := newString("")
	err := conversion.ConvertFrom(str, "boom")
	if err != nil {
		t.Error(err)
	}
	if *str != "boom" {
		t.Errorf("failed to convert boom")
	}
	err = conversion.ConvertFrom(str, newString("bam"))
	if err != nil {
		t.Error(err)
	}
	if *str != "bam" {
		t.Errorf("failed to convert bam")
	}
}

func TestBool(t *testing.T) {
	b := newBool(true)
	err := conversion.ConvertFrom(b, false)
	if err != nil {
		t.Error(err)
	}
	if *b != false {
		t.Errorf("failed to convert false")
	}
	*b = true
	err = conversion.ConvertFrom(b, newBool(false))
	if err != nil {
		t.Error(err)
	}
	if *b != false {
		t.Errorf("failed to convert false")
	}
}

func TestFloat(t *testing.T) {
	f := newFloat32(0)
	err := conversion.ConvertFrom(f, float32(1.2))
	if err != nil {
		t.Error(err)
	}
	if *f < 1.1 || *f > 1.3 {
		t.Errorf("failed to convert 1.2")
	}
	*f = 0.0
	err = conversion.ConvertFrom(f, float64(1.2))
	if err != nil {
		t.Error(err)
	}
	if *f < 1.1 || *f > 1.3 {
		t.Errorf("failed to convert 1.2")
	}
	d := newFloat32(0)
	err = conversion.ConvertFrom(d, float32(1.2))
	if err != nil {
		t.Error(err)
	}
	if *d < 1.1 || *d > 1.3 {
		t.Errorf("failed to convert 1.2")
	}
	*d = 0.0
	err = conversion.ConvertFrom(d, float64(1.2))
	if err != nil {
		t.Error(err)
	}
	if *d < 1.1 || *d > 1.3 {
		t.Errorf("failed to convert 1.2")
	}
}

func TestSlice(t *testing.T) {
	from := make([]int32, 5)
	to := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	}
	err := conversion.ConvertFrom(&from, to)
	if err != nil {
		t.Error(err)
	}
	if len(from) != len(to) {
		t.Errorf("from size: %d, %d", len(from), len(to))
	}
}

func TestMap(t *testing.T) {
	from := make(map[int8]int16, 5)
	to := map[int32]int64{
		1: 1, 2: 2, 3: 3, 4: 4, 5: 5, 6: 6, 7: 7, 8: 8, 9: 9, 10: 10,
	}
	err := conversion.ConvertFrom(&from, to)
	if err != nil {
		t.Error(err)
	}
	if len(from) != len(to) {
		t.Errorf("from size: %d, %d", len(from), len(to))
	}
}

type StructFrom struct {
	A int8
	B int16
	C string
	D float32
	E bool
}

type StructTo struct {
	D float64
	e bool
	B int32
	C string
	F bool
}

func TestStruct(t *testing.T) {
	from := &StructFrom{
		A: 22,
	}
	to := StructTo{
		D: 1.2,
		e: true,
		B: 15,
		C: "boom",
		F: true,
	}
	err := conversion.ConvertFrom(from, to)
	if err != nil {
		t.Fatal(err)
	}
	if from.A != 22 {
		t.Errorf("A should not move")
	}
	if int32(from.B) != to.B {
		t.Errorf("B not the same: %#v", from)
	}
	if from.C != to.C {
		t.Errorf("C not the same: %#v", from)
	}
	if float64(from.D) < to.D-0.1 || float64(from.D) > to.D+0.1 {
		t.Errorf("D not in range: %#v", from)
	}
	if from.E != to.e {
		t.Errorf("E not in range: %#v", from)
	}
}

func TestStructSlice(t *testing.T) {
	from := make([]StructFrom, 2)
	to := []StructTo{
		StructTo{
			D: 1.2,
			e: true,
			B: 15,
			C: "boom",
			F: true,
		}, StructTo{
			D: 1.2,
			e: true,
			B: 15,
			C: "boom",
			F: true,
		},
	}

	err := conversion.ConvertFrom(&from, to)
	if err != nil {
		t.Fatal(err)
	}
}
