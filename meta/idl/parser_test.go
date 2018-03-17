package idl_test

import (
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/type/object"
	"reflect"
	"strings"
	"testing"
)

func helpParserTest(t *testing.T, label, idlSignature string, metaObj *object.MetaObject) {
	reader := strings.NewReader(idlSignature)
	metaObj, err := idl.Parse(reader)
	if err != nil {
		t.Errorf("%s: failed to parse idl", label)
	}
	if !reflect.DeepEqual(metaObj, object.MetaService0) {
		t.Errorf("%s: expected %#v, got %#v", label, object.MetaService0, metaObj)
	}
}

func TestParseServiceServer(t *testing.T) {
	var idl string = `interface Server
	fn authenticate(P0: Map<str,any>) -> Map<str,any>
end
`
	helpParserTest(t, "Service 0", idl, &object.MetaService0)
}
