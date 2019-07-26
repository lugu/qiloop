package object_test

import (
	"strings"
	"testing"

	"github.com/lugu/qiloop/type/object"
)

func TestMetaObjectDecorator(t *testing.T) {
	service0 := object.MetaService0
	id, err := service0.MethodID("authenticate")
	if err != nil {
		panic(err)
	}
	if id != object.AuthenticateActionID {
		t.Errorf("not expecting: %d", id)
	}
	_, err = service0.MethodID("unknown")
	if err == nil {
		panic("shall fail")
	}

	name, err := service0.ActionName(object.AuthenticateActionID)
	if err != nil {
		panic(err)
	}
	if name != "authenticate" {
		panic("invalid name " + name)
	}

	obj := object.FullMetaObject(service0)
	id, err = obj.SignalID("traceObject")
	if err != nil {
		panic(err)
	}
	if id != 0x56 {
		panic("unexpected id")
	}
	_, err = obj.SignalID("unknown")
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
		if signalName != strings.Title(s.Name) {
			t.Errorf("incoherent name: %s and %s", signalName,
				strings.Title(s.Name))
		}
		return nil
	}
	property := func(p object.MetaProperty, propertyName string) error {
		if propertyName != strings.Title(p.Name) {
			t.Errorf("incoherent name: %s and %s", propertyName,
				strings.Title(p.Name))
		}
		return nil
	}
	obj.ForEachMethodAndSignal(method, signal, property)
}

func TestMetaObjectJson(t *testing.T) {
	if object.ObjectMetaObject.JSON() == "" {
		t.Errorf("not expecting empty json")
	}
}
