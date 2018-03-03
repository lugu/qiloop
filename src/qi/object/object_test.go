package object_test

import (
	"bytes"
	"qi/object"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func helpReadGolden(t *testing.T) object.MetaObject {
	path := filepath.Join("testdata", "metaObject-reply-data.bin")
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
