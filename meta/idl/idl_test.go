package idl

import (
	"github.com/lugu/qiloop/type/object"
	"strings"
	"testing"
)

func TestServiceServer(t *testing.T) {
	var w strings.Builder
	if err := GenerateIDL(&w, "Server", object.MetaService0); err != nil {
		t.Errorf("failed to parse server: %s", err)
	}
	expected := `interface Server
	fn authenticate(P0 Map<str,Value>) -> Map<str,Value>
end
`
	if w.String() != expected {
		t.Errorf("Got:\n%s\nExpecting:\n%s\n", w.String(), expected)
	}
}

func TestObject(t *testing.T) {
	var w strings.Builder
	if err := GenerateIDL(&w, "Server", object.ObjectMetaObject); err != nil {
		t.Errorf("failed to parse server: %s", err)
	}
	expected := `interface Server
	fn registerEvent(P0 uint, P1 uint, P2 ulong) -> ulong
	fn unregisterEvent(P0 uint, P1 uint, P2 ulong)
	fn metaObject(P0 uint) -> MetaObject
	fn terminate(P0 uint)
	fn property(P0 Value) -> Value
	fn setProperty(P0 Value, P1 Value)
	fn properties() -> Vec<str>
	fn registerEventWithSignature(P0 uint, P1 uint, P2 ulong, P3 str) -> ulong
end
`
	if w.String() != expected {
		t.Errorf("Got:\n%s\nExpecting:\n%s\n", w.String(), expected)
	}
}
