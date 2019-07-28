// file generated. DO NOT EDIT.
package object

import (
	"fmt"
	basic "github.com/lugu/qiloop/type/basic"
	"io"
)

// MetaMethodParameter is serializable
type MetaMethodParameter struct {
	Name        string
	Description string
}

// readMetaMethodParameter unmarshalls MetaMethodParameter
func readMetaMethodParameter(r io.Reader) (s MetaMethodParameter, err error) {
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Description field: %s", err)
	}
	return s, nil
}

// writeMetaMethodParameter marshalls MetaMethodParameter
func writeMetaMethodParameter(s MetaMethodParameter, w io.Writer) (err error) {
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("write Description field: %s", err)
	}
	return nil
}

// MetaMethod is serializable
type MetaMethod struct {
	Uid                 uint32
	ReturnSignature     string
	Name                string
	ParametersSignature string
	Description         string
	Parameters          []MetaMethodParameter
	ReturnDescription   string
}

// readMetaMethod unmarshalls MetaMethod
func readMetaMethod(r io.Reader) (s MetaMethod, err error) {
	if s.Uid, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read Uid field: %s", err)
	}
	if s.ReturnSignature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read ReturnSignature field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	if s.ParametersSignature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read ParametersSignature field: %s", err)
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Description field: %s", err)
	}
	if s.Parameters, err = func() (b []MetaMethodParameter, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("read slice size: %s", err)
		}
		b = make([]MetaMethodParameter, size)
		for i := 0; i < int(size); i++ {
			b[i], err = readMetaMethodParameter(r)
			if err != nil {
				return b, fmt.Errorf("read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("read Parameters field: %s", err)
	}
	if s.ReturnDescription, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read ReturnDescription field: %s", err)
	}
	return s, nil
}

// writeMetaMethod marshalls MetaMethod
func writeMetaMethod(s MetaMethod, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("write Uid field: %s", err)
	}
	if err := basic.WriteString(s.ReturnSignature, w); err != nil {
		return fmt.Errorf("write ReturnSignature field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	if err := basic.WriteString(s.ParametersSignature, w); err != nil {
		return fmt.Errorf("write ParametersSignature field: %s", err)
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("write Description field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Parameters)), w)
		if err != nil {
			return fmt.Errorf("write slice size: %s", err)
		}
		for _, v := range s.Parameters {
			err = writeMetaMethodParameter(v, w)
			if err != nil {
				return fmt.Errorf("write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("write Parameters field: %s", err)
	}
	if err := basic.WriteString(s.ReturnDescription, w); err != nil {
		return fmt.Errorf("write ReturnDescription field: %s", err)
	}
	return nil
}

// MetaSignal is serializable
type MetaSignal struct {
	Uid       uint32
	Name      string
	Signature string
}

// readMetaSignal unmarshalls MetaSignal
func readMetaSignal(r io.Reader) (s MetaSignal, err error) {
	if s.Uid, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read Uid field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	if s.Signature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Signature field: %s", err)
	}
	return s, nil
}

// writeMetaSignal marshalls MetaSignal
func writeMetaSignal(s MetaSignal, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("write Uid field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	if err := basic.WriteString(s.Signature, w); err != nil {
		return fmt.Errorf("write Signature field: %s", err)
	}
	return nil
}

// MetaProperty is serializable
type MetaProperty struct {
	Uid       uint32
	Name      string
	Signature string
}

// readMetaProperty unmarshalls MetaProperty
func readMetaProperty(r io.Reader) (s MetaProperty, err error) {
	if s.Uid, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read Uid field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Name field: %s", err)
	}
	if s.Signature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Signature field: %s", err)
	}
	return s, nil
}

// writeMetaProperty marshalls MetaProperty
func writeMetaProperty(s MetaProperty, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("write Uid field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("write Name field: %s", err)
	}
	if err := basic.WriteString(s.Signature, w); err != nil {
		return fmt.Errorf("write Signature field: %s", err)
	}
	return nil
}

// MetaObject is serializable
type MetaObject struct {
	Methods     map[uint32]MetaMethod
	Signals     map[uint32]MetaSignal
	Properties  map[uint32]MetaProperty
	Description string
}

// readMetaObject unmarshalls MetaObject
func readMetaObject(r io.Reader) (s MetaObject, err error) {
	if s.Methods, err = func() (m map[uint32]MetaMethod, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("read map size: %s", err)
		}
		m = make(map[uint32]MetaMethod, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(r)
			if err != nil {
				return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
			}
			v, err := readMetaMethod(r)
			if err != nil {
				return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("read Methods field: %s", err)
	}
	if s.Signals, err = func() (m map[uint32]MetaSignal, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("read map size: %s", err)
		}
		m = make(map[uint32]MetaSignal, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(r)
			if err != nil {
				return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
			}
			v, err := readMetaSignal(r)
			if err != nil {
				return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("read Signals field: %s", err)
	}
	if s.Properties, err = func() (m map[uint32]MetaProperty, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("read map size: %s", err)
		}
		m = make(map[uint32]MetaProperty, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(r)
			if err != nil {
				return m, fmt.Errorf("read map key (%d/%d): %s", i+1, size, err)
			}
			v, err := readMetaProperty(r)
			if err != nil {
				return m, fmt.Errorf("read map value (%d/%d): %s", i+1, size, err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("read Properties field: %s", err)
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read Description field: %s", err)
	}
	return s, nil
}

// writeMetaObject marshalls MetaObject
func writeMetaObject(s MetaObject, w io.Writer) (err error) {
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Methods)), w)
		if err != nil {
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range s.Methods {
			err = basic.WriteUint32(k, w)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = writeMetaMethod(v, w)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("write Methods field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Signals)), w)
		if err != nil {
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range s.Signals {
			err = basic.WriteUint32(k, w)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = writeMetaSignal(v, w)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("write Signals field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Properties)), w)
		if err != nil {
			return fmt.Errorf("write map size: %s", err)
		}
		for k, v := range s.Properties {
			err = basic.WriteUint32(k, w)
			if err != nil {
				return fmt.Errorf("write map key: %s", err)
			}
			err = writeMetaProperty(v, w)
			if err != nil {
				return fmt.Errorf("write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("write Properties field: %s", err)
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("write Description field: %s", err)
	}
	return nil
}

// ObjectReference is serializable
type ObjectReference struct {
	MetaObject MetaObject
	ServiceID  uint32
	ObjectID   uint32
}

// readObjectReference unmarshalls ObjectReference
func readObjectReference(r io.Reader) (s ObjectReference, err error) {
	if s.MetaObject, err = readMetaObject(r); err != nil {
		return s, fmt.Errorf("read MetaObject field: %s", err)
	}
	if s.ServiceID, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read ServiceID field: %s", err)
	}
	if s.ObjectID, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read ObjectID field: %s", err)
	}
	return s, nil
}

// writeObjectReference marshalls ObjectReference
func writeObjectReference(s ObjectReference, w io.Writer) (err error) {
	if err := writeMetaObject(s.MetaObject, w); err != nil {
		return fmt.Errorf("write MetaObject field: %s", err)
	}
	if err := basic.WriteUint32(s.ServiceID, w); err != nil {
		return fmt.Errorf("write ServiceID field: %s", err)
	}
	if err := basic.WriteUint32(s.ObjectID, w); err != nil {
		return fmt.Errorf("write ObjectID field: %s", err)
	}
	return nil
}
