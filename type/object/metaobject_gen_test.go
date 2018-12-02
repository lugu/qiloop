package object_test

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/type/object"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func helpReadGolden(t *testing.T) object.MetaObject {
	path := filepath.Join("testdata", "meta-object.bin")
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	metaObj, err := object.ReadMetaObject(file)
	if err != nil {
		t.Errorf("failed to read MetaObject: %s", err)
	}
	return metaObj
}

func TestReadMetaObject(t *testing.T) {
	helpReadGolden(t)
}

func TestReadWriteMetaObject(t *testing.T) {
	metaObj := helpReadGolden(t)
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := object.WriteMetaObject(metaObj, buf); err != nil {
		t.Errorf("failed to write MetaObject: %s", err)
	}
	if metaObj2, err := object.ReadMetaObject(buf); err != nil {
		t.Errorf("failed to re-read MetaObject: %s", err)
	} else if !reflect.DeepEqual(metaObj, metaObj2) {
		t.Errorf("expected %#v, got %#v", metaObj, metaObj2)
	}
}

func TestWriteRead(t *testing.T) {
	var buf bytes.Buffer
	err := object.WriteMetaObject(object.ObjectMetaObject, &buf)
	if err != nil {
		panic(err)
	}
	_, err = object.ReadMetaObject(&buf)
	if err != nil {
		panic(err)
	}
}

func LimitedMetaReader(meta object.MetaObject, size int) io.Reader {
	var buf bytes.Buffer
	object.WriteMetaObject(meta, &buf)
	return &io.LimitedReader{
		R: &buf,
		N: int64(size),
	}
}

func LimitedRefReader(ref object.ObjectReference, size int) io.Reader {
	var buf bytes.Buffer
	object.WriteObjectReference(ref, &buf)
	return &io.LimitedReader{
		R: &buf,
		N: int64(size),
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

func newMetaObject() object.MetaObject {
	return object.MetaObject{
		Description: "Object",
		Methods: map[uint32]object.MetaMethod{
			0x0: {
				Uid:                 0x0,
				ReturnSignature:     "L",
				Name:                "registerEvent",
				ParametersSignature: "(IIL)",
				Parameters: []object.MetaMethodParameter{
					object.MetaMethodParameter{
						Name:        "id",
						Description: "an id",
					},
					object.MetaMethodParameter{
						Name:        "id2",
						Description: "another id",
					},
					object.MetaMethodParameter{
						Name:        "longvalue",
						Description: "another id",
					},
				},
				ReturnDescription: "an id",
			},
			0x1: {
				Uid:                 0x1,
				ReturnSignature:     "v",
				Name:                "unregisterEvent",
				ParametersSignature: "(IIL)",
			},
		},
		Signals: map[uint32]object.MetaSignal{
			0x3: {
				Uid:       0x3,
				Name:      "traceObject",
				Signature: "(is)",
			},
		},
		Properties: map[uint32]object.MetaProperty{
			0x4: {
				Uid:       0x4,
				Name:      "prop",
				Signature: "(s)",
			},
		},
	}
}

func TestWriterMetaObjectError(t *testing.T) {
	meta := newMetaObject()
	var buf bytes.Buffer
	err := object.WriteMetaObject(meta, &buf)
	if err != nil {
		panic(err)
	}
	max := len(buf.Bytes())

	for i := 0; i < max-1; i++ {
		w := NewLimitedWriter(i)
		err := object.WriteMetaObject(meta, w)
		if err == nil {
			panic(fmt.Errorf("not expecting a success at %d", i))
		}
	}
	w := NewLimitedWriter(max)
	err = object.WriteMetaObject(meta, w)
	if err != nil {
		panic(err)
	}
}

func TestReadMetaObjectError(t *testing.T) {
	meta := newMetaObject()
	var buf bytes.Buffer
	err := object.WriteMetaObject(meta, &buf)
	if err != nil {
		panic(err)
	}
	max := len(buf.Bytes())

	for i := 0; i < max; i++ {
		r := LimitedMetaReader(meta, i)
		_, err := object.ReadMetaObject(r)
		if err == nil {
			panic(fmt.Errorf("not expecting a success at %d", i))
		}
	}
	r := LimitedMetaReader(meta, max)
	_, err = object.ReadMetaObject(r)
	if err != nil {
		panic(err)
	}
}

func newObjectReference() object.ObjectReference {
	return object.ObjectReference{
		Boolean:    true,
		MetaObject: newMetaObject(),
		ParentID:   1,
		ServiceID:  1,
		ObjectID:   2,
	}
}

func TestWriterObjectReferenceError(t *testing.T) {
	ref := newObjectReference()
	var buf bytes.Buffer
	err := object.WriteObjectReference(ref, &buf)
	if err != nil {
		panic(err)
	}
	max := len(buf.Bytes())

	for i := 0; i < max-1; i++ {
		w := NewLimitedWriter(i)
		err := object.WriteObjectReference(ref, w)
		if err == nil {
			panic(fmt.Errorf("not expecting a success at %d", i))
		}
	}
	w := NewLimitedWriter(max)
	err = object.WriteObjectReference(ref, w)
	if err != nil {
		panic(err)
	}
}

func TestReadObjectReferenceError(t *testing.T) {
	ref := newObjectReference()
	var buf bytes.Buffer
	err := object.WriteObjectReference(ref, &buf)
	if err != nil {
		panic(err)
	}
	max := len(buf.Bytes())

	for i := 0; i < max; i++ {
		r := LimitedRefReader(ref, i)
		_, err := object.ReadObjectReference(r)
		if err == nil {
			panic(fmt.Errorf("not expecting a success at %d", i))
		}
	}
	r := LimitedRefReader(ref, max)
	_, err = object.ReadObjectReference(r)
	if err != nil {
		panic(err)
	}
}
