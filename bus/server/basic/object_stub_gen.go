// file generated. DO NOT EDIT.
package basic

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	net "github.com/lugu/qiloop/bus/net"
	server "github.com/lugu/qiloop/bus/server"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
	"io"
)

type Object interface {
	Activate(sess bus.Session, serviceID, objectID uint32, signal ObjectSignalHelper) error
	OnTerminate()
	RegisterEvent(P0 uint32, P1 uint32, P2 uint64) (uint64, error)
	UnregisterEvent(P0 uint32, P1 uint32, P2 uint64) error
	MetaObject(P0 uint32) (MetaObject, error)
	Terminate(P0 uint32) error
	Property(P0 value.Value) (value.Value, error)
	SetProperty(P0 value.Value, P1 value.Value) error
	Properties() ([]string, error)
	RegisterEventWithSignature(P0 uint32, P1 uint32, P2 uint64, P3 string) (uint64, error)
}
type ObjectSignalHelper interface {
	SignalTraceObject(P0 EventTrace) error
}
type stubObject struct {
	obj  *server.BasicObject
	impl Object
}

func ObjectObject(impl Object) server.Object {
	var stb stubObject
	stb.impl = impl
	stb.obj = server.NewObject(stb.metaObject())
	stb.obj.Wrap(uint32(0x0), stb.RegisterEvent)
	stb.obj.Wrap(uint32(0x1), stb.UnregisterEvent)
	stb.obj.Wrap(uint32(0x2), stb.MetaObject)
	stb.obj.Wrap(uint32(0x3), stb.Terminate)
	stb.obj.Wrap(uint32(0x5), stb.Property)
	stb.obj.Wrap(uint32(0x6), stb.SetProperty)
	stb.obj.Wrap(uint32(0x7), stb.Properties)
	stb.obj.Wrap(uint32(0x8), stb.RegisterEventWithSignature)
	return &stb
}
func (s *stubObject) Activate(sess bus.Session, serviceID, objectID uint32) error {
	s.obj.Activate(sess, serviceID, objectID)
	return s.impl.Activate(sess, serviceID, objectID, s)
}
func (s *stubObject) OnTerminate() {
	s.impl.OnTerminate()
	s.obj.OnTerminate()
}
func (s *stubObject) Receive(msg *net.Message, from *server.Context) error {
	return s.obj.Receive(msg, from)
}
func (s *stubObject) RegisterEvent(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	P0, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P0: %s", err)
	}
	P1, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P1: %s", err)
	}
	P2, err := basic.ReadUint64(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P2: %s", err)
	}
	ret, callErr := s.impl.RegisterEvent(P0, P1, P2)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteUint64(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (s *stubObject) UnregisterEvent(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	P0, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P0: %s", err)
	}
	P1, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P1: %s", err)
	}
	P2, err := basic.ReadUint64(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P2: %s", err)
	}
	callErr := s.impl.UnregisterEvent(P0, P1, P2)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (s *stubObject) MetaObject(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	P0, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P0: %s", err)
	}
	ret, callErr := s.impl.MetaObject(P0)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := WriteMetaObject(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (s *stubObject) Terminate(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	P0, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P0: %s", err)
	}
	callErr := s.impl.Terminate(P0)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (s *stubObject) Property(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	P0, err := value.NewValue(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P0: %s", err)
	}
	ret, callErr := s.impl.Property(P0)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := ret.Write(&out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (s *stubObject) SetProperty(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	P0, err := value.NewValue(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P0: %s", err)
	}
	P1, err := value.NewValue(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P1: %s", err)
	}
	callErr := s.impl.SetProperty(P0, P1)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	return out.Bytes(), nil
}
func (s *stubObject) Properties(payload []byte) ([]byte, error) {
	ret, callErr := s.impl.Properties()
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := func() error {
		err := basic.WriteUint32(uint32(len(ret)), &out)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range ret {
			err = basic.WriteString(v, &out)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}()
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (s *stubObject) RegisterEventWithSignature(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	P0, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P0: %s", err)
	}
	P1, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P1: %s", err)
	}
	P2, err := basic.ReadUint64(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P2: %s", err)
	}
	P3, err := basic.ReadString(buf)
	if err != nil {
		return nil, fmt.Errorf("cannot read P3: %s", err)
	}
	ret, callErr := s.impl.RegisterEventWithSignature(P0, P1, P2, P3)
	if callErr != nil {
		return nil, callErr
	}
	var out bytes.Buffer
	errOut := basic.WriteUint64(ret, &out)
	if errOut != nil {
		return nil, fmt.Errorf("cannot write response: %s", errOut)
	}
	return out.Bytes(), nil
}
func (s *stubObject) SignalTraceObject(P0 EventTrace) error {
	var buf bytes.Buffer
	if err := WriteEventTrace(P0, &buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	err := s.obj.UpdateSignal(uint32(0x56), buf.Bytes())

	if err != nil {
		return fmt.Errorf("failed to update SignalTraceObject: %s", err)
	}
	return nil
}
func (s *stubObject) metaObject() object.MetaObject {
	return object.MetaObject{
		Description: "Object",
		Methods: map[uint32]object.MetaMethod{
			uint32(0x0): {
				Name:                "registerEvent",
				ParametersSignature: "(IIL)",
				ReturnSignature:     "L",
				Uid:                 uint32(0x0),
			},
			uint32(0x1): {
				Name:                "unregisterEvent",
				ParametersSignature: "(IIL)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x1),
			},
			uint32(0x2): {
				Name:                "metaObject",
				ParametersSignature: "(I)",
				ReturnSignature:     "({I(Issss[(ss)<MetaMethodParameter,name,description>]s)<MetaMethod,uid,returnSignature,name,parametersSignature,description,parameters,returnDescription>}{I(Iss)<MetaSignal,uid,name,signature>}{I(Iss)<MetaProperty,uid,name,signature>}s)<MetaObject,methods,signals,properties,description>",
				Uid:                 uint32(0x2),
			},
			uint32(0x3): {
				Name:                "terminate",
				ParametersSignature: "(I)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x3),
			},
			uint32(0x5): {
				Name:                "property",
				ParametersSignature: "(m)",
				ReturnSignature:     "m",
				Uid:                 uint32(0x5),
			},
			uint32(0x6): {
				Name:                "setProperty",
				ParametersSignature: "(mm)",
				ReturnSignature:     "v",
				Uid:                 uint32(0x6),
			},
			uint32(0x7): {
				Name:                "properties",
				ParametersSignature: "()",
				ReturnSignature:     "[s]",
				Uid:                 uint32(0x7),
			},
			uint32(0x8): {
				Name:                "registerEventWithSignature",
				ParametersSignature: "(IILs)",
				ReturnSignature:     "L",
				Uid:                 uint32(0x8),
			},
		},
		Signals: map[uint32]object.MetaSignal{uint32(0x56): {
			Name:      "traceObject",
			Signature: "((IiIm(ll)<timeval,tv_sec,tv_usec>llII)<EventTrace,id,kind,slotId,arguments,timestamp,userUsTime,systemUsTime,callerContext,calleeContext>)",
			Uid:       uint32(0x56),
		}},
	}
}

type timeval struct {
	Tv_sec  int64
	Tv_usec int64
}

func Readtimeval(r io.Reader) (s timeval, err error) {
	if s.Tv_sec, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read Tv_sec field: " + err.Error())
	}
	if s.Tv_usec, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read Tv_usec field: " + err.Error())
	}
	return s, nil
}
func Writetimeval(s timeval, w io.Writer) (err error) {
	if err := basic.WriteInt64(s.Tv_sec, w); err != nil {
		return fmt.Errorf("failed to write Tv_sec field: " + err.Error())
	}
	if err := basic.WriteInt64(s.Tv_usec, w); err != nil {
		return fmt.Errorf("failed to write Tv_usec field: " + err.Error())
	}
	return nil
}

type EventTrace struct {
	Id            uint32
	Kind          int32
	SlotId        uint32
	Arguments     value.Value
	Timestamp     timeval
	UserUsTime    int64
	SystemUsTime  int64
	CallerContext uint32
	CalleeContext uint32
}

func ReadEventTrace(r io.Reader) (s EventTrace, err error) {
	if s.Id, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read Id field: " + err.Error())
	}
	if s.Kind, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("failed to read Kind field: " + err.Error())
	}
	if s.SlotId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read SlotId field: " + err.Error())
	}
	if s.Arguments, err = value.NewValue(r); err != nil {
		return s, fmt.Errorf("failed to read Arguments field: " + err.Error())
	}
	if s.Timestamp, err = Readtimeval(r); err != nil {
		return s, fmt.Errorf("failed to read Timestamp field: " + err.Error())
	}
	if s.UserUsTime, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read UserUsTime field: " + err.Error())
	}
	if s.SystemUsTime, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read SystemUsTime field: " + err.Error())
	}
	if s.CallerContext, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read CallerContext field: " + err.Error())
	}
	if s.CalleeContext, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read CalleeContext field: " + err.Error())
	}
	return s, nil
}
func WriteEventTrace(s EventTrace, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Id, w); err != nil {
		return fmt.Errorf("failed to write Id field: " + err.Error())
	}
	if err := basic.WriteInt32(s.Kind, w); err != nil {
		return fmt.Errorf("failed to write Kind field: " + err.Error())
	}
	if err := basic.WriteUint32(s.SlotId, w); err != nil {
		return fmt.Errorf("failed to write SlotId field: " + err.Error())
	}
	if err := s.Arguments.Write(w); err != nil {
		return fmt.Errorf("failed to write Arguments field: " + err.Error())
	}
	if err := Writetimeval(s.Timestamp, w); err != nil {
		return fmt.Errorf("failed to write Timestamp field: " + err.Error())
	}
	if err := basic.WriteInt64(s.UserUsTime, w); err != nil {
		return fmt.Errorf("failed to write UserUsTime field: " + err.Error())
	}
	if err := basic.WriteInt64(s.SystemUsTime, w); err != nil {
		return fmt.Errorf("failed to write SystemUsTime field: " + err.Error())
	}
	if err := basic.WriteUint32(s.CallerContext, w); err != nil {
		return fmt.Errorf("failed to write CallerContext field: " + err.Error())
	}
	if err := basic.WriteUint32(s.CalleeContext, w); err != nil {
		return fmt.Errorf("failed to write CalleeContext field: " + err.Error())
	}
	return nil
}

type MetaMethodParameter struct {
	Name        string
	Description string
}

func ReadMetaMethodParameter(r io.Reader) (s MetaMethodParameter, err error) {
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: " + err.Error())
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Description field: " + err.Error())
	}
	return s, nil
}
func WriteMetaMethodParameter(s MetaMethodParameter, w io.Writer) (err error) {
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: " + err.Error())
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("failed to write Description field: " + err.Error())
	}
	return nil
}

type MetaMethod struct {
	Uid                 uint32
	ReturnSignature     string
	Name                string
	ParametersSignature string
	Description         string
	Parameters          []MetaMethodParameter
	ReturnDescription   string
}

func ReadMetaMethod(r io.Reader) (s MetaMethod, err error) {
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
			b[i], err = ReadMetaMethodParameter(r)
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
func WriteMetaMethod(s MetaMethod, w io.Writer) (err error) {
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
			err = WriteMetaMethodParameter(v, w)
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

type MetaSignal struct {
	Uid       uint32
	Name      string
	Signature string
}

func ReadMetaSignal(r io.Reader) (s MetaSignal, err error) {
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
func WriteMetaSignal(s MetaSignal, w io.Writer) (err error) {
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

type MetaProperty struct {
	Uid       uint32
	Name      string
	Signature string
}

func ReadMetaProperty(r io.Reader) (s MetaProperty, err error) {
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
func WriteMetaProperty(s MetaProperty, w io.Writer) (err error) {
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

type MetaObject struct {
	Methods     map[uint32]MetaMethod
	Signals     map[uint32]MetaSignal
	Properties  map[uint32]MetaProperty
	Description string
}

func ReadMetaObject(r io.Reader) (s MetaObject, err error) {
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
			v, err := ReadMetaMethod(r)
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
			v, err := ReadMetaSignal(r)
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
			v, err := ReadMetaProperty(r)
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
func WriteMetaObject(s MetaObject, w io.Writer) (err error) {
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
			err = WriteMetaMethod(v, w)
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
			err = WriteMetaSignal(v, w)
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
			err = WriteMetaProperty(v, w)
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
