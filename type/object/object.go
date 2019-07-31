package object

import (
	"github.com/lugu/qiloop/type/value"
)

// MetaObjectMethodID is the action id of the method MetaObject.
const MetaObjectMethodID = 2

// MinUserActionID is the first custom action id.
const MinUserActionID = 100

// Object represents an object on QiMessaging. Every services
// implement the Object interface except the Server service which only
// has one method (authenticate).
type Object interface {

	// MetaObject returns a description of an Object.
	MetaObject(objectID uint32) (MetaObject, error)

	// Terminate informs a service that an object is no longer needed.
	Terminate(objectID uint32) error

	// Subscribes to a signal or a property. Return an handler or
	// an error.
	RegisterEvent(objectID uint32, actionID uint32, handler uint64) (
		uint64, error)
	// Unsubscribe to a signal or a property.
	UnregisterEvent(objectID uint32, actionID uint32, handler uint64) error

	// RegisterEventWithSignature: is similar with RegisterEvent:
	// it subscribes to a signal. The value of the signal is
	// converted to the type described by the signature.
	// Not supported.
	RegisterEventWithSignature(objectID uint32, signalID uint32,
		handler uint64, signature string) (uint64, error)

	// Property returns the value of the property.
	Property(prop value.Value) (value.Value, error)

	// SetProperty update the value of a property.
	SetProperty(prop value.Value, val value.Value) error

	// Properties returns the list of properties.
	Properties() ([]string, error)
}
