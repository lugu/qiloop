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
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Description field: " + err.Error())
	}
	return s, nil
}

// writeMetaMethodParameter marshalls MetaMethodParameter
func writeMetaMethodParameter(s MetaMethodParameter, w io.Writer) (err error) {
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("failed to write Description field: " + err.Error())
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
		return s, fmt.Errorf("failed to read Uid field: " + err.Error())
	}
	if s.ReturnSignature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read ReturnSignature field: " + err.Error())
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	if s.ParametersSignature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read ParametersSignature field: " + err.Error())
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Description field: " + err.Error())
	}
	if s.Parameters, err = func() (b []MetaMethodParameter, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]MetaMethodParameter, size)
		for i := 0; i < int(size); i++ {
			b[i], err = readMetaMethodParameter(r)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Parameters field: " + err.Error())
	}
	if s.ReturnDescription, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read ReturnDescription field: " + err.Error())
	}
	return s, nil
}

// writeMetaMethod marshalls MetaMethod
func writeMetaMethod(s MetaMethod, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("failed to write Uid field: " + err.Error())
	}
	if err := basic.WriteString(s.ReturnSignature, w); err != nil {
		return fmt.Errorf("failed to write ReturnSignature field: " + err.Error())
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
	}
	if err := basic.WriteString(s.ParametersSignature, w); err != nil {
		return fmt.Errorf("failed to write ParametersSignature field: " + err.Error())
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("failed to write Description field: " + err.Error())
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Parameters)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.Parameters {
			err = writeMetaMethodParameter(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Parameters field: " + err.Error())
	}
	if err := basic.WriteString(s.ReturnDescription, w); err != nil {
		return fmt.Errorf("failed to write ReturnDescription field: " + err.Error())
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
		return s, fmt.Errorf("failed to read Uid field: " + err.Error())
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	if s.Signature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Signature field: " + err.Error())
	}
	return s, nil
}

// writeMetaSignal marshalls MetaSignal
func writeMetaSignal(s MetaSignal, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("failed to write Uid field: " + err.Error())
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
	}
	if err := basic.WriteString(s.Signature, w); err != nil {
		return fmt.Errorf("failed to write Signature field: " + err.Error())
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
		return s, fmt.Errorf("failed to read Uid field: " + err.Error())
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	if s.Signature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Signature field: " + err.Error())
	}
	return s, nil
}

// writeMetaProperty marshalls MetaProperty
func writeMetaProperty(s MetaProperty, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("failed to write Uid field: " + err.Error())
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
	}
	if err := basic.WriteString(s.Signature, w); err != nil {
		return fmt.Errorf("failed to write Signature field: " + err.Error())
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
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[uint32]MetaMethod, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := readMetaMethod(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Methods field: " + err.Error())
	}
	if s.Signals, err = func() (m map[uint32]MetaSignal, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[uint32]MetaSignal, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := readMetaSignal(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Signals field: " + err.Error())
	}
	if s.Properties, err = func() (m map[uint32]MetaProperty, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[uint32]MetaProperty, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := readMetaProperty(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Properties field: " + err.Error())
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Description field: " + err.Error())
	}
	return s, nil
}

// writeMetaObject marshalls MetaObject
func writeMetaObject(s MetaObject, w io.Writer) (err error) {
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Methods)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.Methods {
			err = basic.WriteUint32(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = writeMetaMethod(v, w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Methods field: " + err.Error())
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Signals)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.Signals {
			err = basic.WriteUint32(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = writeMetaSignal(v, w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Signals field: " + err.Error())
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Properties)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.Properties {
			err = basic.WriteUint32(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = writeMetaProperty(v, w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Properties field: " + err.Error())
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("failed to write Description field: " + err.Error())
	}
	return nil
}

// ObjectReference is serializable
type ObjectReference struct {
	Boolean    bool
	MetaObject MetaObject
	ParentID   uint32
	ServiceID  uint32
	ObjectID   uint32
}

// readObjectReference unmarshalls ObjectReference
func readObjectReference(r io.Reader) (s ObjectReference, err error) {
	if s.Boolean, err = basic.ReadBool(r); err != nil {
		return s, fmt.Errorf("failed to read Boolean field: " + err.Error())
	}
	if s.MetaObject, err = readMetaObject(r); err != nil {
		return s, fmt.Errorf("failed to read MetaObject field: " + err.Error())
	}
	if s.ParentID, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read ParentID field: " + err.Error())
	}
	if s.ServiceID, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read ServiceID field: " + err.Error())
	}
	if s.ObjectID, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read ObjectID field: " + err.Error())
	}
	return s, nil
}

// writeObjectReference marshalls ObjectReference
func writeObjectReference(s ObjectReference, w io.Writer) (err error) {
	if err := basic.WriteBool(s.Boolean, w); err != nil {
		return fmt.Errorf("failed to write Boolean field: " + err.Error())
	}
	if err := writeMetaObject(s.MetaObject, w); err != nil {
		return fmt.Errorf("failed to write MetaObject field: " + err.Error())
	}
	if err := basic.WriteUint32(s.ParentID, w); err != nil {
		return fmt.Errorf("failed to write ParentID field: " + err.Error())
	}
	if err := basic.WriteUint32(s.ServiceID, w); err != nil {
		return fmt.Errorf("failed to write ServiceID field: " + err.Error())
	}
	if err := basic.WriteUint32(s.ObjectID, w); err != nil {
		return fmt.Errorf("failed to write ObjectID field: " + err.Error())
	}
	return nil
}
