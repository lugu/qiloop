package object_test

import (
	"github.com/lugu/qiloop/type/object"
	"strings"
	"testing"
)

func TestMetaObjectDecorator(t *testing.T) {
	service0 := &object.MetaService0
	id, err := service0.MethodUid("authenticate")
	if err != nil {
		panic(err)
	}
	if id != object.AuthenticateActionID {
		t.Errorf("not expecting: %d", id)
	}
	_, err = service0.MethodUid("unknown")
	if err == nil {
		panic("shall fail")
	}

	obj := object.FullMetaObject(*service0)
	id, err = obj.SignalUid("traceObject")
	if err != nil {
		panic(err)
	}
	if id != 0x56 {
		panic("unexpected id")
	}
	_, err = obj.SignalUid("unknown")
	if err == nil {
		panic("shall fail")
	}
	method := func(m object.MetaMethod, methodName string) error {
		if methodName != strings.Title(m.Name) {
			t.Errorf("incoherent name: %s and %s", methodName,
				strings.Title(m.Name))
		}
		return nil
	}
	signal := func(s object.MetaSignal, signalName string) error {
		if signalName != "Signal"+strings.Title(s.Name) {
			t.Errorf("incoherent name: %s and %s", signalName,
				strings.Title(s.Name))
		}
		return nil
	}
	obj.ForEachMethodAndSignal(method, signal)
}

func TestMetaObjectJson(t *testing.T) {
	if object.ObjectMetaObject.Json() == "" {
		t.Errorf("not expecting empty json")
	}
}
