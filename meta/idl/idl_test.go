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
	fn authenticate
end
`
	if w.String() != expected {
		t.Errorf("Got:\n%s\nExpecting:\n%s\n", w.String(), expected)
	}
}
