package encoding_test

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/lugu/qiloop/type/encoding"
)

type Info struct {
	A int32
	B int16
	C int64
}

func (i *Info) Encode(e encoding.Encoder) error {
	if err := e.Encode(i.A); err != nil {
		return fmt.Errorf("field a: %w", err)
	}
	if err := e.Encode(i.B); err != nil {
		return fmt.Errorf("field b: %w", err)
	}
	if err := e.Encode(i.C); err != nil {
		return fmt.Errorf("field c: %w", err)
	}
	return nil
}

// pointer since it is modifying the struct
func (i *Info) Decode(d encoding.Decoder) error {
	if err := d.Decode(&i.A); err != nil {
		return fmt.Errorf("field a: %w", err)
	}
	if err := d.Decode(&i.B); err != nil {
		return fmt.Errorf("field b: %w", err)
	}
	if err := d.Decode(&i.C); err != nil {
		return fmt.Errorf("field c: %w", err)
	}
	return nil
}

func TestBasic(t *testing.T) {
	var buf bytes.Buffer
	permission := make(map[string]string)

	encoder := encoding.NewEncoder(permission, &buf)

	if err := encoder.Encode(int16(1)); err != nil {
		t.Error(err)
	}
	if err := encoder.Encode(int32(2)); err != nil {
		t.Error(err)
	}
	if err := encoder.Encode(int64(3)); err != nil {
		t.Error(err)
	}

	decoder := encoding.NewDecoder(permission, &buf)

	var a int16
	var b int32
	var c int64

	if err := decoder.Decode(&a); err != nil {
		t.Error(err)
	}
	if a != 1 {
		t.Errorf("unexpected value: %d", a)
	}
	if err := decoder.Decode(&b); err != nil {
		t.Error(err)
	}
	if b != 2 {
		t.Errorf("unexpected value: %d", a)
	}
	if err := decoder.Decode(&c); err != nil {
		t.Error(err)
	}
	if c != 3 {
		t.Errorf("unexpected value: %d", a)
	}
}

func TestStruct(t *testing.T) {
	var buf bytes.Buffer
	permission := make(map[string]string)

	in := &Info{1, 2, 3}
	encoder := encoding.NewEncoder(permission, &buf)
	if err := encoder.Encode(in); err != nil {
		t.Error(err)
	}

	var out Info
	decoder := encoding.NewDecoder(permission, &buf)
	if err := decoder.Decode(&out); err != nil {
		t.Error(err)
	}
	if out.A != 1 || out.B != 2 || out.C != 3 {
		t.Errorf("unexpected value: %v", out)
	}
}

func TestDecodeSlice(t *testing.T) {
	b := bytes.NewBuffer([]byte{
		0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0x0,
	})
	permission := make(map[string]string)
	decoder := encoding.NewDecoder(permission, b)
	var a []int16
	err := decoder.Decode(&a)
	if err != nil {
		t.Error(err)
	}
	if a[0] != 1 || a[1] != 2 {
		t.Errorf("Unexpected value: %v", a)
	}
}

func newInt16(n int16) *int16 { return &n }
func newInt32(n int32) *int32 { return &n }
func newInt64(n int64) *int64 { return &n }

var testData = map[string]struct {
	wire []byte
	in   interface{}
	out  interface{}
}{
	"int16": {[]byte{
		0x01, 0x00,
	}, int16(1), newInt16(1)},
	"int32": {[]byte{
		0x02, 0x00, 0x00, 0x00,
	}, int32(2), newInt32(2)},
	"int64": {[]byte{
		0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, int64(3), newInt64(3)},
	"struct": {[]byte{
		0x01, 0x00, 0x00, 0x00,
		0x02, 0x00,
		0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, Info{1, 2, 3}, &Info{1, 2, 3}},
	"list of int16": {[]byte{
		0x02, 0x00, 0x00, 0x00,
		0x01, 0x00,
		0x02, 0x00,
	}, []int16{1, 2}, &[]int16{1, 2}},
	"list of struct": {[]byte{
		0x01, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00,
		0x02, 0x00,
		0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, []Info{Info{1, 2, 3}}, &[]Info{Info{1, 2, 3}}},
	"map of int16": {[]byte{
		0x02, 0x00, 0x00, 0x00,
		0x01, 0x00,
		0x02, 0x00,
		0x03, 0x00,
		0x04, 0x00,
	}, map[int16]int16{1: 2, 3: 4}, &map[int16]int16{1: 2, 3: 4}},
	"map of struct": {[]byte{
		0x02, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00,
		0x02, 0x00,
		0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x00, 0x00,
		0x04, 0x00,
		0x05, 0x00,
		0x06, 0x00,
		0x07, 0x00,
		0x08, 0x00, 0x00, 0x00,
		0x09, 0x00,
		0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x00, 0x00,
		0x0b, 0x00,
		0x0c, 0x00,
		0x0d, 0x00,
		0x0e, 0x00,
	}, map[Info][]int16{
		Info{1, 2, 3}:  []int16{4, 5, 6, 7},
		Info{8, 9, 10}: []int16{11, 12, 13, 14},
	}, &map[Info][]int16{
		Info{1, 2, 3}:  []int16{4, 5, 6, 7},
		Info{8, 9, 10}: []int16{11, 12, 13, 14},
	}},
}

func TestDecode(t *testing.T) {
	permission := make(map[string]string)
	for name, test := range testData {
		buf := bytes.NewBuffer(test.wire)
		decoder := encoding.NewDecoder(permission, buf)
		pv := reflect.New(reflect.TypeOf(test.out).Elem())
		val := pv.Interface()
		err := decoder.Decode(val)
		if err != nil {
			t.Errorf("Decode failed for %s: %v", name, err)
		}
		if !reflect.DeepEqual(val, test.out) {
			t.Errorf("%s:\nhave %#v\nwant %#v", name, val, test.out)
		}
	}
}

func TestEncode(t *testing.T) {
	permission := make(map[string]string)
	for name, test := range testData {
		var buf bytes.Buffer
		encoder := encoding.NewEncoder(permission, &buf)
		err := encoder.Encode(test.in)
		if err != nil {
			t.Errorf("Encode failed for %s %v", name, err)
		}
		decoder := encoding.NewDecoder(permission, &buf)
		pv := reflect.New(reflect.TypeOf(test.out).Elem())
		val := pv.Interface()
		err = decoder.Decode(val)
		if err != nil {
			t.Errorf("Decode failed for %s: %v", name, err)
		}
		if !reflect.DeepEqual(val, test.out) {
			t.Errorf("%s:\nhave %#v\nwant %#v", name, val, test.in)
		}
	}
}

type InfoReflect struct {
	A int32
	B int16
	C int64
}

func helpBenchPrepareDecoder(n int) io.Reader {
	infoBytes := []byte{
		0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x03, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00,
	}
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		buf.Write(infoBytes)
	}
	return &buf
}

func BenchmarkReadInfoBaseLine(b *testing.B) {
	reader := helpBenchPrepareDecoder(b.N)
	a := make([]InfoBaseLine, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		a[i], err = readInfoBaseLine(reader)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		if a[i].A != 1 {
			b.Errorf("wrong A at index: %d", i)
		}
		if a[i].B != 2 {
			b.Errorf("wrong B at index: %d", i)
		}
		if a[i].C != 3 {
			b.Errorf("wrong C at index: %d", i)
		}
	}
}

func BenchmarkReadInfoCodeGen(b *testing.B) {
	reader := helpBenchPrepareDecoder(b.N)
	permission := make(map[string]string)
	decoder := encoding.NewDecoder(permission, reader)
	a := make([]Info, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := decoder.Decode(&a[i])
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		if a[i].A != 1 {
			b.Errorf("wrong A at index: %d", i)
		}
		if a[i].B != 2 {
			b.Errorf("wrong B at index: %d", i)
		}
		if a[i].C != 3 {
			b.Errorf("wrong C at index: %d", i)
		}
	}
}

func BenchmarkReadInfoReflect(b *testing.B) {
	reader := helpBenchPrepareDecoder(b.N)
	permission := make(map[string]string)
	decoder := encoding.NewDecoder(permission, reader)
	a := make([]InfoReflect, b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := decoder.Decode(&a[i])
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		if a[i].A != 1 {
			b.Errorf("wrong A at index: %d", i)
		}
		if a[i].B != 2 {
			b.Errorf("wrong B at index: %d", i)
		}
		if a[i].C != 3 {
			b.Errorf("wrong C at index: %d", i)
		}
	}
}

func BenchmarkReadInfoGob(b *testing.B) {
	var buf bytes.Buffer
	encoder := encoding.NewGobEncoder(&buf)
	a := make([]Info, b.N)
	for i := 0; i < b.N; i++ {
		a[i] = Info{1, 2, 3}
	}
	err := encoder.Encode(a)
	if err != nil {
		b.Error(err)
	}
	decoder := encoding.NewGobDecoder(&buf)
	a = make([]Info, b.N)
	b.ResetTimer()
	err = decoder.Decode(&a)
	if err != nil {
		b.Error(err)
	}
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		if a[i].A != 1 {
			b.Errorf("wrong A at index: %d", i)
		}
		if a[i].B != 2 {
			b.Errorf("wrong B at index: %d", i)
		}
		if a[i].C != 3 {
			b.Errorf("wrong C at index: %d", i)
		}
	}
}

func BenchmarkReadInfoJSON(b *testing.B) {
	var buf bytes.Buffer
	encoder := encoding.NewJSONEncoder(&buf)
	a := make([]Info, b.N)
	for i := 0; i < b.N; i++ {
		a[i] = Info{1, 2, 3}
	}
	err := encoder.Encode(a)
	if err != nil {
		b.Error(err)
	}
	decoder := encoding.NewJSONDecoder(buf.Bytes())
	a = make([]Info, b.N)
	b.ResetTimer()
	err = decoder.Decode(&a)
	if err != nil {
		b.Error(err)
	}
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		if a[i].A != 1 {
			b.Errorf("wrong A at index: %d", i)
		}
		if a[i].B != 2 {
			b.Errorf("wrong B at index: %d", i)
		}
		if a[i].C != 3 {
			b.Errorf("wrong C at index: %d", i)
		}
	}
}

func BenchmarkWriteInfoBaseLine(b *testing.B) {
	writer := bytes.NewBuffer(make([]byte, b.N*8))
	a := InfoBaseLine{1, 2, 3}
	for i := 0; i < b.N; i++ {
		var err error
		err = writeInfoBaseLine(a, writer)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkWriteInfoCodeGen(b *testing.B) {
	writer := bytes.NewBuffer(make([]byte, b.N*8))
	permission := make(map[string]string)
	encoder := encoding.NewEncoder(permission, writer)
	a := Info{1, 2, 3}
	for i := 0; i < b.N; i++ {
		var err error
		err = encoder.Encode(a)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkWriteInfoReflect(b *testing.B) {
	writer := bytes.NewBuffer(make([]byte, b.N*8))
	permission := make(map[string]string)
	encoder := encoding.NewEncoder(permission, writer)
	a := InfoReflect{1, 2, 3}
	for i := 0; i < b.N; i++ {
		var err error
		err = encoder.Encode(a)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkWriteJSON(b *testing.B) {
	writer := bytes.NewBuffer(make([]byte, b.N*8))
	encoder := encoding.NewJSONEncoder(writer)
	a := Info{1, 2, 3}
	for i := 0; i < b.N; i++ {
		var err error
		err = encoder.Encode(a)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkWriteGob(b *testing.B) {
	writer := bytes.NewBuffer(make([]byte, b.N*8))
	encoder := encoding.NewGobEncoder(writer)
	a := Info{1, 2, 3}
	for i := 0; i < b.N; i++ {
		var err error
		err = encoder.Encode(a)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func helpJSONSerialization(t *testing.T, name string, in interface{}) {
	var buf bytes.Buffer
	encoder := encoding.NewJSONEncoder(&buf)
	err := encoder.Encode(in)
	if err != nil {
		t.Errorf("%s: encode error: %w", name, err)
		return
	}
	decoder := encoding.NewJSONDecoder(buf.Bytes())
	pv := reflect.New(reflect.TypeOf(in).Elem())
	val := pv.Interface()
	err = decoder.Decode(val)
	if err != nil {
		t.Errorf("Decode failed for %s: %v", name, err)
	}
	if !reflect.DeepEqual(val, in) {
		t.Errorf("%s: have %#v\nwant %#v", name, in, val)
	}
}

func TestJSONSerialization(t *testing.T) {
	for name, test := range testData {
		if name == "map of struct" {
			// for some reason it cannot be encoded into JSON...
			continue
		}
		helpJSONSerialization(t, name, test.out)
	}
}

func helpGobSerialization(t *testing.T, name string, in interface{}) {
	var buf bytes.Buffer
	encoder := encoding.NewGobEncoder(&buf)
	err := encoder.Encode(in)
	if err != nil {
		t.Errorf("%s: encode error: %w", name, err)
		return
	}
	decoder := encoding.NewGobDecoder(&buf)
	pv := reflect.New(reflect.TypeOf(in).Elem())
	val := pv.Interface()
	err = decoder.Decode(val)
	if err != nil {
		t.Errorf("Decode failed for %s: %v", name, err)
	}
	if !reflect.DeepEqual(val, in) {
		t.Errorf("%s: have %#v\nwant %#v", name, in, val)
	}
}

func TestGobSerialization(t *testing.T) {
	for name, test := range testData {
		helpGobSerialization(t, name, test.out)
	}
}
