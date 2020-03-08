package encoding_test

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/lugu/qiloop/type/encoding"
	"github.com/lugu/qiloop/type/value"
)

type InfoReflect struct {
	A int16
	B int32
	C int64
	D uint16
	E uint32
	F uint64
	G string
	H bool
	I float32
	J float64
}

type Info struct {
	A int16
	B int32
	C int64
	D uint16
	E uint32
	F uint64
	G string
	H bool
	I float32
	J float64
}

var info = Info{1, 2, 3, 1, 2, 3, "hello", true, 1.5, 3.0}
var base = InfoBaseLine{1, 2, 3, 1, 2, 3, "hello", true, 1.5, 3.0}
var ref = InfoReflect{1, 2, 3, 1, 2, 3, "hello", true, 1.5, 3.0}

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
	if err := e.Encode(i.D); err != nil {
		return fmt.Errorf("field d: %w", err)
	}
	if err := e.Encode(i.E); err != nil {
		return fmt.Errorf("field e: %w", err)
	}
	if err := e.Encode(i.F); err != nil {
		return fmt.Errorf("field f: %w", err)
	}
	if err := e.Encode(i.G); err != nil {
		return fmt.Errorf("field g: %w", err)
	}
	if err := e.Encode(i.H); err != nil {
		return fmt.Errorf("field h: %w", err)
	}
	if err := e.Encode(i.I); err != nil {
		return fmt.Errorf("field i: %w", err)
	}
	if err := e.Encode(i.J); err != nil {
		return fmt.Errorf("field j: %w", err)
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
	if err := d.Decode(&i.D); err != nil {
		return fmt.Errorf("field d: %w", err)
	}
	if err := d.Decode(&i.E); err != nil {
		return fmt.Errorf("field e: %w", err)
	}
	if err := d.Decode(&i.F); err != nil {
		return fmt.Errorf("field f: %w", err)
	}
	if err := d.Decode(&i.G); err != nil {
		return fmt.Errorf("field g: %w", err)
	}
	if err := d.Decode(&i.H); err != nil {
		return fmt.Errorf("field h: %w", err)
	}
	if err := d.Decode(&i.I); err != nil {
		return fmt.Errorf("field i: %w", err)
	}
	if err := d.Decode(&i.J); err != nil {
		return fmt.Errorf("field j: %w", err)
	}
	return nil
}

func TestBasic(t *testing.T) {
	var buf bytes.Buffer

	permission := encoding.DefaultCap()
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
	permission := encoding.DefaultCap()

	in := ref
	encoder := encoding.NewEncoder(permission, &buf)
	if err := encoder.Encode(in); err != nil {
		t.Error(err)
	}

	var out Info
	decoder := encoding.NewDecoder(permission, &buf)
	if err := decoder.Decode(&out); err != nil {
		t.Error(err)
	}
	if out.A != in.A || out.B != in.B || out.C != in.C ||
		out.D != in.D || out.E != in.E || out.F != in.F ||
		out.G != in.G || out.H != in.H || out.I != in.I ||
		out.J != in.J {
		t.Errorf("unexpected value: %v", out)
	}
}

func TestDecodeSlice(t *testing.T) {
	b := bytes.NewBuffer([]byte{
		0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0x0,
	})
	permission := encoding.DefaultCap()
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

func newBool(b bool) *bool       { return &b }
func newString(s string) *string { return &s }

func newUint(n uint) *uint       { return &n }
func newUint16(n uint16) *uint16 { return &n }
func newUint32(n uint32) *uint32 { return &n }
func newUint64(n uint64) *uint64 { return &n }

func newInt(n int) *int       { return &n }
func newInt16(n int16) *int16 { return &n }
func newInt32(n int32) *int32 { return &n }
func newInt64(n int64) *int64 { return &n }

var testData = map[string]interface{}{
	"bool-true":      newBool(true),
	"bool-false":     newBool(false),
	"string-1":       newString("A"),
	"string-2":       newString("A"),
	"uint16":         newUint16(1),
	"uint32":         newUint32(2),
	"uint64":         newUint64(3),
	"uint":           newUint(3),
	"int16":          newInt16(1),
	"int32":          newInt32(2),
	"int64":          newInt64(3),
	"int":            newInt(3),
	"struct-1":       &ref,
	"struct-2":       &info,
	"list of int16":  &[]int16{1, 2},
	"list of struct": &[]Info{info},
	"map of int16":   &map[int16]int16{1: 2, 3: 4},
	"map of struct": &map[Info][]int16{
		info: []int16{4, 5, 6, 7},
	},
}

func TestEncode(t *testing.T) {
	permission := encoding.DefaultCap()
	for name, test := range testData {
		var buf bytes.Buffer
		encoder := encoding.NewEncoder(permission, &buf)
		err := encoder.Encode(test)
		if err != nil {
			t.Errorf("Encode failed for %s %v", name, err)
		}
		decoder := encoding.NewDecoder(permission, &buf)
		pv := reflect.New(reflect.TypeOf(test).Elem())
		val := pv.Interface()
		err = decoder.Decode(val)
		if err != nil {
			t.Errorf("Decode failed for %s: %v", name, err)
		}
		if !reflect.DeepEqual(val, test) {
			var buf bytes.Buffer
			encoder := encoding.NewEncoder(permission, &buf)
			err := encoder.Encode(test)
			if err != nil {
				t.Errorf("Encode failed for %s %v", name, err)
			}

			t.Errorf("%s:\nhave %#v\nwant %#v\n%#v",
				name, val, test, buf.Bytes())
		}
	}
}

func TestMapValue(t *testing.T) {
	test := map[string]value.Value{
		"a": value.Bool(true),
		"b": value.Int(12),
		"c": value.String("oups"),
	}
	var buf bytes.Buffer
	permission := encoding.DefaultCap()
	encoder := encoding.NewEncoder(permission, &buf)
	err := encoder.Encode(test)
	if err != nil {
		t.Errorf("Encode failed: %v", err)
	}
	decoder := encoding.NewDecoder(permission, &buf)
	val := map[string]value.Value{}
	err = decoder.Decode(&val)
	if err != nil {
		t.Errorf("Decode failed: %v", err)
	}
	if !reflect.DeepEqual(val, test) {
		t.Errorf("Not the same: %#v vs %#v", val, test)
	}
}

func TestSliceValue(t *testing.T) {
	test := []value.Value{
		value.Bool(true),
		value.Int(12),
		value.String("oups"),
	}
	var buf bytes.Buffer
	permission := encoding.DefaultCap()
	encoder := encoding.NewEncoder(permission, &buf)
	err := encoder.Encode(test)
	if err != nil {
		t.Errorf("Encode failed: %v", err)
	}
	decoder := encoding.NewDecoder(permission, &buf)
	val := []value.Value{}
	err = decoder.Decode(&val)
	if err != nil {
		t.Errorf("Decode failed: %v", err)
	}
	if !reflect.DeepEqual(val, test) {
		t.Errorf("Not the same: %#v vs %#v", val, test)
	}
}

func helpBenchPrepareDecoder(n int) io.Reader {
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		permission := encoding.DefaultCap()
		encoder := encoding.NewEncoder(permission, &buf)
		err := encoder.Encode(ref)
		if err != nil {
			panic(err)
		}
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
	permission := encoding.DefaultCap()
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
	permission := encoding.DefaultCap()
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
		a[i] = info
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
		a[i] = info
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
	a := base
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
	permission := encoding.DefaultCap()
	encoder := encoding.NewEncoder(permission, writer)
	a := info
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
	permission := encoding.DefaultCap()
	encoder := encoding.NewEncoder(permission, writer)
	a := ref
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
	a := info
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
	a := info
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
		helpJSONSerialization(t, name, test)
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
		helpGobSerialization(t, name, test)
	}
}
