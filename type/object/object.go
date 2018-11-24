package object

import (
	"github.com/lugu/qiloop/type/value"
)

const MetaObjectMethodID = 2

const MinUserActionID = 100

// Object represents an object on QiMessaging. Every services
// implement the Object interface except the Server service which only
// has one method (authenticate).
type Object interface {

	// MetaObject returns a description of an Object.
	MetaObject(objectID uint32) (MetaObject, error)

	// Terminate informs a service that an object is no longer needed.
	Terminate(objectID uint32) error

	RegisterEvent(serviceID uint32, signalID uint32, handler uint64) (uint64, error)
	UnregisterEvent(serviceID uint32, signalID uint32, handler uint64) error
	RegisterEventWithSignature(serviceID uint32, signalID uint32, handler uint64, signature string) (uint64, error)

	// Property returns the value of the property.
	Property(prop value.Value) (value.Value, error)

	// SetProperty update the value of a property.
	SetProperty(prop value.Value, val value.Value) error

	// Properties returns the list of properties.
	Properties() ([]string, error)
}
