package idl

import (
	"github.com/lugu/qiloop/meta/signature"
	"testing"
)

func TestScopeSearch(t *testing.T) {
	sc := NewScope()
	err := sc.Add("test", signature.NewLongType())
	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = sc.Search("test")
	if err != nil {
		t.Errorf("%s", err)
	}
	_, err = sc.Search("test-missing")
	if err == nil {
		t.Errorf("an error expected")
	}
}

func TestScopeCollision(t *testing.T) {
	sc := NewScope()
	err := sc.Add("test", signature.NewLongType())
	if err != nil {
		t.Errorf("%s", err)
	}
	err = sc.Add("test", signature.NewIntType())
	if err == nil {
		t.Errorf("an error expected")
	}
}

func TestScopeGlobal(t *testing.T) {
	sc1 := NewScope()
	err := sc1.Add("test", signature.NewLongType())
	if err != nil {
		t.Errorf("%s", err)
	}
	sc0 := NewScope()
	err = sc0.Add("test", signature.NewIntType())
	if err != nil {
		t.Errorf("%s", err)
	}
	err = sc0.Extend("sc1", sc1)
	if err != nil {
		t.Errorf("%s", err)
	}

	typ, err := sc0.Search("test")
	if err != nil {
		t.Errorf("%s", err)
	}
	if typ.Signature() != "i" {
		t.Errorf("Unexpected signature: %s", typ.Signature())
	}
	typ, err = sc0.Search("sc1.test")
	if err != nil {
		t.Errorf("%s", err)
	}
	if typ.Signature() != "l" {
		t.Errorf("Unexpected signature: %s", typ.Signature())
	}
}

func TestScopeGlobalCollision(t *testing.T) {
	sc2 := NewScope()
	sc1 := NewScope()
	sc0 := NewScope()
	err := sc0.Extend("sc1", sc1)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = sc0.Extend("sc1", sc2)
	if err == nil {
		t.Errorf("Error expected")
	}
}
