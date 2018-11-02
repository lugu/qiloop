package basic

import (
	"github.com/lugu/qiloop/bus/session"
)

type BasicObject struct {
	Meta MetaObject
}

func (b *BasicObject) Activate(sess session.Session, serviceID, objectID uint32, signal ObjectSignalHelper) {
	panic("not yet implemented")
}
func (b *BasicObject) RegisterEvent(P0 uint32, P1 uint32, P2 uint64) (uint64, error) {
	panic("not yet implemented")
}
func (b *BasicObject) UnregisterEvent(P0 uint32, P1 uint32, P2 uint64) error {
	panic("not yet implemented")
}
func (b *BasicObject) MetaObject(P0 uint32) (MetaObject, error) {
	panic("not yet implemented")
}
func (b *BasicObject) Terminate(P0 uint32) error {
	panic("not yet implemented")
}
func (b *BasicObject) Property(P0 value.Value) (value.Value, error) {
	panic("not yet implemented")
}
func (b *BasicObject) SetProperty(P0 value.Value, P1 value.Value) error {
	panic("not yet implemented")
}
func (b *BasicObject) Properties() ([]string, error) {
	panic("not yet implemented")
}
func (b *BasicObject) RegisterEventWithSignature(P0 uint32, P1 uint32, P2 uint64, P3 string) (uint64, error) {
	panic("not yet implemented")
}
