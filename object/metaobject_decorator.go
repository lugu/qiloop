package object

import (
	fmt "fmt"
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
