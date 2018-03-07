package stage2_test

import (
	"bytes"
	"github.com/lugu/qiloop/meta/stage2"
	"os"
	"reflect"
	"testing"
)

func helpReadGolden(t *testing.T) stage2.MetaObject {
	file, err := os.Open("meta-stage2.bin")
	if err != nil {
		t.Fatal(err)
	}
	metaObj, err := stage2.ReadMetaObject(file)
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
	if err := stage2.WriteMetaObject(metaObj, buf); err != nil {
		t.Errorf("failed to write MetaObject: %s", err)
	}
	if metaObj2, err := stage2.ReadMetaObject(buf); err != nil {
		t.Errorf("failed to re-read MetaObject: %s", err)
	} else if !reflect.DeepEqual(metaObj, metaObj2) {
		t.Errorf("expected %#v, got %#v", metaObj, metaObj2)
	}
}
