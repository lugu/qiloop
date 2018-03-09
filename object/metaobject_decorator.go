package object

import (
	fmt "fmt"
)

// MethodUid returns the identifier of the method based on its name.
func (m *MetaObject) MethodUid(name string) (uint32, error) {
	for k, method := range m.Methods {
		if method.Name == name {
			return k, nil
		}
	}
	return 0, fmt.Errorf("failed to find method %s", name)
}
