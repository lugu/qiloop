package object

import (
	"fmt"
	"sort"
	"strings"
)

// SignalUid returns the ID of a method given its name.
func (m *MetaObject) MethodUid(name string) (uint32, error) {
	for k, method := range m.Methods {
		if method.Name == name {
			return k, nil
		}
	}
	return 0, fmt.Errorf("failed to find method %s", name)
}

// SignalUid returns the ID of a signal given its name.
func (m *MetaObject) SignalUid(name string) (uint32, error) {
	for k, signal := range m.Signals {
		if signal.Name == name {
			return k, nil
		}
	}
	return 0, fmt.Errorf("failed to find signal %s", name)
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

// ForEachMethodAndSignal call methodCall and signalCall for each
// method and signal respectively. Always call them in the same order
// and unsure the generated method names are unique within the object.
func (metaObj *MetaObject) ForEachMethodAndSignal(
	methodCall func(m MetaMethod, methodName string) error,
	signalCall func(s MetaSignal, methodName string) error) error {

	methodNames := make(map[string]bool)
	keys := make([]int, 0)
	for k := range metaObj.Methods {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, i := range keys {
		k := uint32(i)
		m := metaObj.Methods[k]

		// generate uniq name for the method
		methodName := registerName(strings.Title(m.Name), methodNames)
		if err := methodCall(m, methodName); err != nil {
			return fmt.Errorf("method callback failed for %s: %s", m.Name, err)
		}
	}
	keys = make([]int, 0)
	for k := range metaObj.Signals {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, i := range keys {
		k := uint32(i)
		s := metaObj.Signals[k]
		methodName := registerName("Signal"+strings.Title(s.Name), methodNames)

		if err := signalCall(s, methodName); err != nil {
			return fmt.Errorf("signal callback failed for %s: %s", s.Name, err)
		}
	}
	return nil
}
