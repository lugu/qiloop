package object_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/lugu/qiloop/meta/signature"
	"github.com/lugu/qiloop/type/object"
)

func helpFullObject() object.MetaObject {
	return object.MetaObject{
		Description: "Object",
		Methods: map[uint32]object.MetaMethod{
			0: {
				Name:                "registerEvent",
				ParametersSignature: "(IIL)",
				ReturnSignature:     "L",
				Uid:                 0,
			},
			1: {
				Name:                "unregisterEvent",
				ParametersSignature: "(IIL)",
				ReturnSignature:     "v",
				Uid:                 1,
			},
			2: {
				Name:                "metaObject",
				ParametersSignature: "(I)",
				ReturnSignature:     "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>",
				Uid:                 2,
			},
			3: {
				Name:                "terminate",
				ParametersSignature: "(I)",
				ReturnSignature:     "v",
				Uid:                 3,
			},
			5: {
				Name:                "property",
				ParametersSignature: "(m)",
				ReturnSignature:     "m",
				Uid:                 5,
			},
			6: {
				Name:                "setProperty",
				ParametersSignature: "(mm)",
				ReturnSignature:     "v",
				Uid:                 6,
			},
			7: {
				Name:                "properties",
				ParametersSignature: "()",
				ReturnSignature:     "[s]",
				Uid:                 7,
			},
			8: {
				Name:                "registerEventWithSignature",
				ParametersSignature: "(IILs)",
				ReturnSignature:     "L",
				Uid:                 8,
			},
			80: {
				Name:                "isStatsEnabled",
				ParametersSignature: "()",
				ReturnSignature:     "b",
				Uid:                 80,
			},
			81: {
				Name:                "enableStats",
				ParametersSignature: "(b)",
				ReturnSignature:     "v",
				Uid:                 81,
			},
			82: {
				Name:                "stats",
				ParametersSignature: "()",
				ReturnSignature:     "{I(I(fff)<MinMaxSum,minValue,maxValue,cumulatedValue>(fff)<MinMaxSum,minValue,maxValue,cumulatedValue>(fff)<MinMaxSum,minValue,maxValue,cumulatedValue>)<MethodStatistics,count,wall,user,system>}",
				Uid:                 82,
			},
			83: {
				Name:                "clearStats",
				ParametersSignature: "()",
				ReturnSignature:     "v",
				Uid:                 83,
			},
			84: {
				Name:                "isTraceEnabled",
				ParametersSignature: "()",
				ReturnSignature:     "b",
				Uid:                 84,
			},
			85: {
				Name:                "enableTrace",
				ParametersSignature: "(b)",
				ReturnSignature:     "v",
				Uid:                 85,
			},
		},
		Properties: map[uint32]object.MetaProperty{},
		Signals: map[uint32]object.MetaSignal{86: {
			Name:      "traceObject",
			Signature: "(IiIm(ll)<timeval,tv_sec,tv_usec>llII)<EventTrace,id,kind,slotId,arguments,timestamp,userUsTime,systemUsTime,callerContext,calleeContext>",
			Uid:       86,
		}},
	}
}

func helpReadGolden(t *testing.T) object.MetaObject {
	path := filepath.Join("testdata", "meta-object.bin")
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	metaObj, err := object.ReadMetaObject(file)
	if err != nil {
		t.Errorf("read MetaObject: %s", err)
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
		t.Errorf("write MetaObject: %s", err)
	}
	if metaObj2, err := object.ReadMetaObject(buf); err != nil {
		t.Errorf("re-read MetaObject: %s", err)
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

func helpMethodID(t *testing.T, name, param, ret string, expected uint32) {
	metaObj := helpReadGolden(t)
	id, _, err := metaObj.MethodID(name, param)
	if err != nil {
		t.Errorf("cannot find %s %s %s: %s", name, param, ret, err)
	} else if id != expected {
		t.Errorf("Wrong id: %d", id)
	}
}

func TestMethodID(t *testing.T) {
	helpMethodID(t, "registerEvent", "(IIL)", "L", 0)
	helpMethodID(t, "unregisterEvent", "(IIL)", "v", 1)
	helpMethodID(t, "property", "(m)", "m", 5)
	helpMethodID(t, "metaObject", "(I)", signature.MetaObjectSignature, 2)
}

func helpSignalID(t *testing.T, name, sig string, expected uint32) {
	metaObj := helpFullObject()
	id, err := metaObj.SignalID(name, sig)
	if err != nil {
		t.Errorf("cannot find %s %s: %s", name, sig, err)
	} else if id != expected {
		t.Errorf("Wrong id: %d", id)
	}
}

func TestSignalID(t *testing.T) {
	helpSignalID(t, "traceObject", "(IiIm(ll)<timeval,tv_sec,tv_usec>llII)<EventTrace,id,kind,slotId,arguments,timestamp,userUsTime,systemUsTime,callerContext,calleeContext>", 86)
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
		MetaObject: newMetaObject(),
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
