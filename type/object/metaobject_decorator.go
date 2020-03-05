package object

import (
	"encoding/json"
	"fmt"
	io "io"
	"sort"
	"strings"

	"github.com/lugu/qiloop/meta/signature"
)

// ActionName returns the name of the method, signal or property.
func (m *MetaObject) ActionName(id uint32) (string, error) {
	for k, method := range m.Methods {
		if k == id {
			return method.Name, nil
		}
	}
	for k, signal := range m.Signals {
		if k == id {
			return signal.Name, nil
		}
	}
	for k, property := range m.Properties {
		if k == id {
			return property.Name, nil
		}
	}
	return "", fmt.Errorf("action not found: %d", id)
}

// MethodID returns the ID of a method given its name, the parameters
// signature and the returned value signature.
func (m *MetaObject) MethodID(name, signature string ) (uint32, string, error) {
	for k, method := range m.Methods {
		if method.Name == name &&
			method.ParametersSignature == signature {
			return k, method.ReturnSignature, nil
		}
	}
	for k, method := range m.Methods {
		if method.Name == name {
			return k, method.ReturnSignature, nil
		}
	}
	return 0, "", fmt.Errorf("missing method %s", name)
}

// SignalID returns the ID of a signal given its name and its
// signature. If the signature does not match:
// - if it is not a tuple, try with a tuple
// - else return the first signal with the same name
func (m *MetaObject) SignalID(name, sig string) (uint32, error) {
	for _, signal := range m.Signals {
		if signal.Name == name &&
			signal.Signature == sig {
			return signal.Uid, nil
		}
	}
	t, err := signature.Parse(sig)
	if err != nil {
		return 0, fmt.Errorf("cannot parse signal signature: %s", err)
	}
	_, ok := t.(*signature.TupleType)
	if !ok {
		// try to wrap the type in a tuple
		tuple := signature.NewTupleType([]signature.Type{t})
		return m.SignalID(name, tuple.Signature())
	}
	// last chance: try the first signal with the same name.
	for _, signal := range m.Signals {
		if signal.Name == name {
			return signal.Uid, nil
		}
	}
	return 0, fmt.Errorf("cannot find signal %s", name)
}

// PropertyID returns the ID of a property given its name and its
// signature. If the signature does not match:
// - if it is not a tuple, try with a tuple
// - else return the first property with the same name
func (m *MetaObject) PropertyID(name, sig string) (uint32, error) {
	for _, property := range m.Properties {
		if property.Name == name &&
			property.Signature == sig {
			return property.Uid, nil
		}
	}
	// The signature may or may not include the tuple parameters.
	// If it does not, add it and test again.
	t, err := signature.Parse(sig)
	if err != nil {
		return 0, fmt.Errorf("cannot parse property signature: %s", err)
	}
	_, ok := t.(*signature.TupleType)
	if !ok {
		// try to wrap the type in a tuple
		tuple := signature.NewTupleType([]signature.Type{t})
		return m.PropertyID(name, tuple.Signature())
	}
	// last chance: try the first signal with the same name.
	for _, property := range m.Properties {
		if property.Name == name {
			return property.Uid, nil
		}
	}
	return 0, fmt.Errorf("cannot find property %s", name)
}

func (m *MetaObject) PropertyName(id uint32) (string, error) {
	prop, ok := m.Properties[id]
	if !ok {
		return "", fmt.Errorf("missing property %d", id)
	}
	return prop.Name, nil
}

func registerName(name string, names map[string]bool) string {
	newName := name
	for i := 0; i < 100; i++ {
		if _, ok := names[newName]; !ok {
			break
		}
		newName = fmt.Sprintf("%s_%d", name, i)
	}
	names[newName] = true
	return newName
}

// JSON returns a json string describing the meta object.
func (m *MetaObject) JSON() string {
	json, _ := json.MarshalIndent(m, "", "    ")
	return string(json)
}

// ForEachMethodAndSignal call methodCall and signalCall for each
// method and signal respectively. Always call them in the same order
// and unsure the generated method names are unique within the object.
func (m *MetaObject) ForEachMethodAndSignal(
	methodCall func(m MetaMethod, methodName string) error,
	signalCall func(s MetaSignal, signalName string) error,
	propertyCall func(p MetaProperty, propertyName string) error) error {

	methodNames := make(map[string]bool)
	keys := make([]int, 0)
	for k := range m.Methods {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, i := range keys {
		k := uint32(i)
		m := m.Methods[k]

		// generate uniq name for the method
		methodName := registerName(strings.Title(m.Name), methodNames)
		if err := methodCall(m, methodName); err != nil {
			return fmt.Errorf("method callback failed for %s: %s", m.Name, err)
		}
	}
	keys = make([]int, 0)
	for k := range m.Signals {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, i := range keys {
		k := uint32(i)
		s := m.Signals[k]
		signalName := registerName(strings.Title(s.Name), methodNames)

		if err := signalCall(s, signalName); err != nil {
			return fmt.Errorf("signal callback failed for %s: %s", s.Name, err)
		}
	}
	keys = make([]int, 0)
	for k := range m.Properties {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, i := range keys {
		k := uint32(i)
		p := m.Properties[k]
		propertyName := registerName(strings.Title(p.Name), methodNames)

		if err := propertyCall(p, propertyName); err != nil {
			return fmt.Errorf("property callback failed for %s: %s",
				p.Name, err)
		}
	}
	return nil
}

// ReadMetaObject unserialize a MetaObject and returns it.
func ReadMetaObject(r io.Reader) (s MetaObject, err error) {
	return readMetaObject(r)
}

// WriteMetaObject serialize a MetaObject.
func WriteMetaObject(s MetaObject, w io.Writer) (err error) {
	return writeMetaObject(s, w)
}

// ReadMetaObject serialize an ObjectReference.
func ReadObjectReference(r io.Reader) (s ObjectReference, err error) {
	return readObjectReference(r)
}

// WriteObjectReference unserialize an ObjectReference and returns it.
func WriteObjectReference(s ObjectReference, w io.Writer) (err error) {
	return writeObjectReference(s, w)
}

// FullMetaObject fills the meta object with generic objects methods.
func FullMetaObject(from MetaObject) MetaObject {

	var meta MetaObject
	meta.Description = from.Description
	meta.Methods = make(map[uint32]MetaMethod)
	meta.Signals = make(map[uint32]MetaSignal)
	meta.Properties = make(map[uint32]MetaProperty)
	for id, m := range from.Methods {
		meta.Methods[id] = m
	}
	for id, s := range from.Signals {
		meta.Signals[id] = s
	}
	for id, p := range from.Properties {
		meta.Properties[id] = p
	}

	for i, method := range ObjectMetaObject.Methods {
		meta.Methods[i] = method
	}
	for i, signal := range ObjectMetaObject.Signals {
		meta.Signals[i] = signal
	}
	for i, property := range ObjectMetaObject.Properties {
		meta.Properties[i] = property
	}
	return meta
}
