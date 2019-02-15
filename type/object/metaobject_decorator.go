package object

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// MethodID returns the ID of a method given its name.
func (m *MetaObject) MethodID(name string) (uint32, error) {
	for k, method := range m.Methods {
		if method.Name == name {
			return k, nil
		}
	}
	return 0, fmt.Errorf("failed to find method %s", name)
}

// SignalID returns the ID of a signal
func (m *MetaObject) SignalID(name string) (uint32, error) {
	for _, signal := range m.Signals {
		if signal.Name == name {
			return signal.Uid, nil
		}
	}
	return 0, fmt.Errorf("failed to find signal %s", name)
}

// PropertyID returns the ID of a property
func (m *MetaObject) PropertyID(name string) (uint32, error) {
	for _, property := range m.Properties {
		if property.Name == name {
			return property.Uid, nil
		}
	}
	return 0, fmt.Errorf("failed to find property %s", name)
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
	signalCall func(s MetaSignal, subscribeMethodName string) error,
	propertyCall func(p MetaProperty, getMethodName, setMethodName,
		subscribeMethodName string) error) error {

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
		methodName := registerName("Subscribe"+strings.Title(s.Name), methodNames)

		if err := signalCall(s, methodName); err != nil {
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
		getMethodName := registerName("Get"+strings.Title(p.Name),
			methodNames)
		setMethodName := registerName("Set"+strings.Title(p.Name),
			methodNames)
		subscribeMethodName := registerName("Subscribe"+strings.Title(
			p.Name), methodNames)

		if err := propertyCall(p, getMethodName, setMethodName,
			subscribeMethodName); err != nil {
			return fmt.Errorf("property callback failed for %s: %s",
				p.Name, err)
		}
	}
	return nil
}

// FullMetaObject fills the meta object with generic objects methods.
func FullMetaObject(meta MetaObject) MetaObject {
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
