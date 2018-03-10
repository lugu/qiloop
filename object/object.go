package object

import (
	"github.com/lugu/qiloop/value"
)

const MetaObjectMethodID = 2

// Object represents an object on QiMessaging. Every services
// implement the Object interface except the Server service which only
// has one method (authenticate).
type Object interface {
	// ServiceID returns the related service identifier
	ServiceID() uint32
	// ServiceID returns object identifier within the service
	// namespace.
	ObjectID() uint32

	// MetaObject returns a description of an Object.
	MetaObject(objectID uint32) (MetaObject, error)

	// Terminate informs a service that an object is no longer needed.
	Terminate(objectID uint32) error

	RegisterEvent(fromServiceID uint32, fromObjectID uint32, signalID uint64) (uint64, error)
	RegisterEventWithSignature(fromServiceID uint32, fromObjectID uint32, signalID uint64, signature string) (uint64, error)
	UnregisterEvent(fromServiceID uint32, fromObjectID uint32, signalID uint64) error

	// Property returns the value of the property.
	Property(prop value.Value) (value.Value, error)

	// SetProperty update the value of a property.
	SetProperty(prop value.Value, val value.Value) error

	// Properties returns the list of properties.
	Properties() ([]string, error)
}
