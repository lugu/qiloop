// file generated. DO NOT EDIT.
package services

import (
	bytes "bytes"
	fmt "fmt"
	basic "github.com/lugu/qiloop/basic"
	object "github.com/lugu/qiloop/object"
	session "github.com/lugu/qiloop/session"
	value "github.com/lugu/qiloop/value"
	io "io"
)

type Server struct {
	session.Proxy
}

func NewServer(ses session.Session, obj uint32) (*Server, error) {
	proxy, err := ses.Proxy("Server", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &Server{proxy}, nil
}
func (p *Server) Authenticate(p0 map[string]value.Value) (map[string]value.Value, error) {
	var err error
	var ret map[string]value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = func() error {
		err := basic.WriteUint32(uint32(len(p0)), buf)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range p0 {
			err = basic.WriteString(k, buf)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = v.Write(buf)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("authenticate", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call authenticate failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (m map[string]value.Value, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]value.Value, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := value.NewValue(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse authenticate response: %s", err)
	}
	return ret, nil
}

type Object struct {
	session.Proxy
}

func NewObject(ses session.Session, obj uint32) (*Object, error) {
	proxy, err := ses.Proxy("Object", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &Object{proxy}, nil
}
func (p *Object) RegisterEvent(p0 uint32, p1 uint32, p2 uint64) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteUint32(p1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p1: %s", err)
	}
	if err = basic.WriteUint64(p2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p2: %s", err)
	}
	response, err := p.Call("registerEvent", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerEvent failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse registerEvent response: %s", err)
	}
	return ret, nil
}
func (p *Object) UnregisterEvent(p0 uint32, p1 uint32, p2 uint64) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteUint32(p1, buf); err != nil {
		return fmt.Errorf("failed to serialize p1: %s", err)
	}
	if err = basic.WriteUint64(p2, buf); err != nil {
		return fmt.Errorf("failed to serialize p2: %s", err)
	}
	_, err = p.Call("unregisterEvent", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterEvent failed: %s", err)
	}
	return nil
}
func (p *Object) MetaObject(p0 uint32) (object.MetaObject, error) {
	var err error
	var ret object.MetaObject
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("metaObject", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call metaObject failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = object.ReadMetaObject(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse metaObject response: %s", err)
	}
	return ret, nil
}
func (p *Object) Terminate(p0 uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("terminate", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call terminate failed: %s", err)
	}
	return nil
}
func (p *Object) Property(p0 value.Value) (value.Value, error) {
	var err error
	var ret value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = p0.Write(buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("property", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call property failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = value.NewValue(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse property response: %s", err)
	}
	return ret, nil
}
func (p *Object) SetProperty(p0 value.Value, p1 value.Value) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = p0.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = p1.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize p1: %s", err)
	}
	_, err = p.Call("setProperty", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setProperty failed: %s", err)
	}
	return nil
}
func (p *Object) Properties() ([]string, error) {
	var err error
	var ret []string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("properties", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call properties failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(buf)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse properties response: %s", err)
	}
	return ret, nil
}
func (p *Object) RegisterEventWithSignature(p0 uint32, p1 uint32, p2 uint64, p3 string) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteUint32(p1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p1: %s", err)
	}
	if err = basic.WriteUint64(p2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p2: %s", err)
	}
	if err = basic.WriteString(p3, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p3: %s", err)
	}
	response, err := p.Call("registerEventWithSignature", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerEventWithSignature failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse registerEventWithSignature response: %s", err)
	}
	return ret, nil
}

type ServiceDirectory struct {
	session.Proxy
}

func NewServiceDirectory(ses session.Session, obj uint32) (*ServiceDirectory, error) {
	proxy, err := ses.Proxy("ServiceDirectory", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &ServiceDirectory{proxy}, nil
}
func (p *ServiceDirectory) RegisterEvent(p0 uint32, p1 uint32, p2 uint64) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteUint32(p1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p1: %s", err)
	}
	if err = basic.WriteUint64(p2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p2: %s", err)
	}
	response, err := p.Call("registerEvent", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerEvent failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse registerEvent response: %s", err)
	}
	return ret, nil
}
func (p *ServiceDirectory) UnregisterEvent(p0 uint32, p1 uint32, p2 uint64) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteUint32(p1, buf); err != nil {
		return fmt.Errorf("failed to serialize p1: %s", err)
	}
	if err = basic.WriteUint64(p2, buf); err != nil {
		return fmt.Errorf("failed to serialize p2: %s", err)
	}
	_, err = p.Call("unregisterEvent", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterEvent failed: %s", err)
	}
	return nil
}
func (p *ServiceDirectory) MetaObject(p0 uint32) (object.MetaObject, error) {
	var err error
	var ret object.MetaObject
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("metaObject", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call metaObject failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = object.ReadMetaObject(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse metaObject response: %s", err)
	}
	return ret, nil
}
func (p *ServiceDirectory) Terminate(p0 uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("terminate", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call terminate failed: %s", err)
	}
	return nil
}
func (p *ServiceDirectory) Property(p0 value.Value) (value.Value, error) {
	var err error
	var ret value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = p0.Write(buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("property", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call property failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = value.NewValue(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse property response: %s", err)
	}
	return ret, nil
}
func (p *ServiceDirectory) SetProperty(p0 value.Value, p1 value.Value) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = p0.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = p1.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize p1: %s", err)
	}
	_, err = p.Call("setProperty", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setProperty failed: %s", err)
	}
	return nil
}
func (p *ServiceDirectory) Properties() ([]string, error) {
	var err error
	var ret []string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("properties", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call properties failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(buf)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse properties response: %s", err)
	}
	return ret, nil
}
func (p *ServiceDirectory) RegisterEventWithSignature(p0 uint32, p1 uint32, p2 uint64, p3 string) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteUint32(p1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p1: %s", err)
	}
	if err = basic.WriteUint64(p2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p2: %s", err)
	}
	if err = basic.WriteString(p3, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p3: %s", err)
	}
	response, err := p.Call("registerEventWithSignature", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerEventWithSignature failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse registerEventWithSignature response: %s", err)
	}
	return ret, nil
}
func (p *ServiceDirectory) IsStatsEnabled() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("isStatsEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isStatsEnabled failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse isStatsEnabled response: %s", err)
	}
	return ret, nil
}
func (p *ServiceDirectory) EnableStats(p0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("enableStats", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableStats failed: %s", err)
	}
	return nil
}
func (p *ServiceDirectory) Stats() (map[uint32]MethodStatistics, error) {
	var err error
	var ret map[uint32]MethodStatistics
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("stats", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call stats failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (m map[uint32]MethodStatistics, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[uint32]MethodStatistics, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := ReadMethodStatistics(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse stats response: %s", err)
	}
	return ret, nil
}
func (p *ServiceDirectory) ClearStats() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("clearStats", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearStats failed: %s", err)
	}
	return nil
}
func (p *ServiceDirectory) IsTraceEnabled() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("isTraceEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isTraceEnabled failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse isTraceEnabled response: %s", err)
	}
	return ret, nil
}
func (p *ServiceDirectory) EnableTrace(p0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("enableTrace", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableTrace failed: %s", err)
	}
	return nil
}
func (p *ServiceDirectory) Service(p0 string) (ServiceInfo, error) {
	var err error
	var ret ServiceInfo
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("service", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call service failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = ReadServiceInfo(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse service response: %s", err)
	}
	return ret, nil
}
func (p *ServiceDirectory) Services() ([]ServiceInfo, error) {
	var err error
	var ret []ServiceInfo
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("services", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call services failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (b []ServiceInfo, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]ServiceInfo, size)
		for i := 0; i < int(size); i++ {
			b[i], err = ReadServiceInfo(buf)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse services response: %s", err)
	}
	return ret, nil
}
func (p *ServiceDirectory) RegisterService(p0 ServiceInfo) (uint32, error) {
	var err error
	var ret uint32
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = WriteServiceInfo(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("registerService", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerService failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadUint32(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse registerService response: %s", err)
	}
	return ret, nil
}
func (p *ServiceDirectory) UnregisterService(p0 uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("unregisterService", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterService failed: %s", err)
	}
	return nil
}
func (p *ServiceDirectory) ServiceReady(p0 uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("serviceReady", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call serviceReady failed: %s", err)
	}
	return nil
}
func (p *ServiceDirectory) UpdateServiceInfo(p0 ServiceInfo) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = WriteServiceInfo(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("updateServiceInfo", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call updateServiceInfo failed: %s", err)
	}
	return nil
}
func (p *ServiceDirectory) MachineId() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("machineId", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call machineId failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse machineId response: %s", err)
	}
	return ret, nil
}
func (p *ServiceDirectory) _socketOfService(p0 uint32) (object.ObjectReference, error) {
	var err error
	var ret object.ObjectReference
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("_socketOfService", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _socketOfService failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = object.ReadObjectReference(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _socketOfService response: %s", err)
	}
	return ret, nil
}

type LogManager struct {
	session.Proxy
}

func NewLogManager(ses session.Session, obj uint32) (*LogManager, error) {
	proxy, err := ses.Proxy("LogManager", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &LogManager{proxy}, nil
}
func (p *LogManager) RegisterEvent(p0 uint32, p1 uint32, p2 uint64) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteUint32(p1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p1: %s", err)
	}
	if err = basic.WriteUint64(p2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p2: %s", err)
	}
	response, err := p.Call("registerEvent", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerEvent failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse registerEvent response: %s", err)
	}
	return ret, nil
}
func (p *LogManager) UnregisterEvent(p0 uint32, p1 uint32, p2 uint64) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteUint32(p1, buf); err != nil {
		return fmt.Errorf("failed to serialize p1: %s", err)
	}
	if err = basic.WriteUint64(p2, buf); err != nil {
		return fmt.Errorf("failed to serialize p2: %s", err)
	}
	_, err = p.Call("unregisterEvent", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterEvent failed: %s", err)
	}
	return nil
}
func (p *LogManager) MetaObject(p0 uint32) (object.MetaObject, error) {
	var err error
	var ret object.MetaObject
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("metaObject", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call metaObject failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = object.ReadMetaObject(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse metaObject response: %s", err)
	}
	return ret, nil
}
func (p *LogManager) Terminate(p0 uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("terminate", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call terminate failed: %s", err)
	}
	return nil
}
func (p *LogManager) Property(p0 value.Value) (value.Value, error) {
	var err error
	var ret value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = p0.Write(buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("property", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call property failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = value.NewValue(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse property response: %s", err)
	}
	return ret, nil
}
func (p *LogManager) SetProperty(p0 value.Value, p1 value.Value) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = p0.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = p1.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize p1: %s", err)
	}
	_, err = p.Call("setProperty", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setProperty failed: %s", err)
	}
	return nil
}
func (p *LogManager) Properties() ([]string, error) {
	var err error
	var ret []string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("properties", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call properties failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(buf)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse properties response: %s", err)
	}
	return ret, nil
}
func (p *LogManager) RegisterEventWithSignature(p0 uint32, p1 uint32, p2 uint64, p3 string) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteUint32(p1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p1: %s", err)
	}
	if err = basic.WriteUint64(p2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p2: %s", err)
	}
	if err = basic.WriteString(p3, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p3: %s", err)
	}
	response, err := p.Call("registerEventWithSignature", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call registerEventWithSignature failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse registerEventWithSignature response: %s", err)
	}
	return ret, nil
}
func (p *LogManager) IsStatsEnabled() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("isStatsEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isStatsEnabled failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse isStatsEnabled response: %s", err)
	}
	return ret, nil
}
func (p *LogManager) EnableStats(p0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("enableStats", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableStats failed: %s", err)
	}
	return nil
}
func (p *LogManager) Stats() (map[uint32]MethodStatistics, error) {
	var err error
	var ret map[uint32]MethodStatistics
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("stats", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call stats failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (m map[uint32]MethodStatistics, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[uint32]MethodStatistics, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadUint32(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := ReadMethodStatistics(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse stats response: %s", err)
	}
	return ret, nil
}
func (p *LogManager) ClearStats() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("clearStats", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearStats failed: %s", err)
	}
	return nil
}
func (p *LogManager) IsTraceEnabled() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("isTraceEnabled", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isTraceEnabled failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse isTraceEnabled response: %s", err)
	}
	return ret, nil
}
func (p *LogManager) EnableTrace(p0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("enableTrace", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableTrace failed: %s", err)
	}
	return nil
}
func (p *LogManager) Log(p0 []LogMessage) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = func() error {
		err := basic.WriteUint32(uint32(len(p0)), buf)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range p0 {
			err = WriteLogMessage(v, buf)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("log", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call log failed: %s", err)
	}
	return nil
}
func (p *LogManager) CreateListener() (object.ObjectReference, error) {
	var err error
	var ret object.ObjectReference
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("createListener", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call createListener failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = object.ReadObjectReference(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse createListener response: %s", err)
	}
	return ret, nil
}
func (p *LogManager) GetListener() (object.ObjectReference, error) {
	var err error
	var ret object.ObjectReference
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("getListener", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getListener failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = object.ReadObjectReference(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse getListener response: %s", err)
	}
	return ret, nil
}
func (p *LogManager) AddProvider(p0 object.ObjectReference) (uint32, error) {
	var err error
	var ret uint32
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = object.WriteObjectReference(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("addProvider", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call addProvider failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadUint32(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse addProvider response: %s", err)
	}
	return ret, nil
}
func (p *LogManager) RemoveProvider(p0 uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("removeProvider", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call removeProvider failed: %s", err)
	}
	return nil
}

type MinMaxSum struct {
	MinValue       float32
	MaxValue       float32
	CumulatedValue float32
}

func ReadMinMaxSum(r io.Reader) (s MinMaxSum, err error) {
	if s.MinValue, err = basic.ReadFloat32(r); err != nil {
		return s, fmt.Errorf("failed to read MinValue field: %s", err)
	}
	if s.MaxValue, err = basic.ReadFloat32(r); err != nil {
		return s, fmt.Errorf("failed to read MaxValue field: %s", err)
	}
	if s.CumulatedValue, err = basic.ReadFloat32(r); err != nil {
		return s, fmt.Errorf("failed to read CumulatedValue field: %s", err)
	}
	return s, nil
}
func WriteMinMaxSum(s MinMaxSum, w io.Writer) (err error) {
	if err := basic.WriteFloat32(s.MinValue, w); err != nil {
		return fmt.Errorf("failed to write MinValue field: %s", err)
	}
	if err := basic.WriteFloat32(s.MaxValue, w); err != nil {
		return fmt.Errorf("failed to write MaxValue field: %s", err)
	}
	if err := basic.WriteFloat32(s.CumulatedValue, w); err != nil {
		return fmt.Errorf("failed to write CumulatedValue field: %s", err)
	}
	return nil
}

type MethodStatistics struct {
	Count  uint32
	Wall   MinMaxSum
	User   MinMaxSum
	System MinMaxSum
}

func ReadMethodStatistics(r io.Reader) (s MethodStatistics, err error) {
	if s.Count, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read Count field: %s", err)
	}
	if s.Wall, err = ReadMinMaxSum(r); err != nil {
		return s, fmt.Errorf("failed to read Wall field: %s", err)
	}
	if s.User, err = ReadMinMaxSum(r); err != nil {
		return s, fmt.Errorf("failed to read User field: %s", err)
	}
	if s.System, err = ReadMinMaxSum(r); err != nil {
		return s, fmt.Errorf("failed to read System field: %s", err)
	}
	return s, nil
}
func WriteMethodStatistics(s MethodStatistics, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Count, w); err != nil {
		return fmt.Errorf("failed to write Count field: %s", err)
	}
	if err := WriteMinMaxSum(s.Wall, w); err != nil {
		return fmt.Errorf("failed to write Wall field: %s", err)
	}
	if err := WriteMinMaxSum(s.User, w); err != nil {
		return fmt.Errorf("failed to write User field: %s", err)
	}
	if err := WriteMinMaxSum(s.System, w); err != nil {
		return fmt.Errorf("failed to write System field: %s", err)
	}
	return nil
}

type ServiceInfo struct {
	Name      string
	ServiceId uint32
	MachineId string
	ProcessId uint32
	Endpoints []string
	SessionId string
}

func ReadServiceInfo(r io.Reader) (s ServiceInfo, err error) {
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: %s", err)
	}
	if s.ServiceId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read ServiceId field: %s", err)
	}
	if s.MachineId, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read MachineId field: %s", err)
	}
	if s.ProcessId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read ProcessId field: %s", err)
	}
	if s.Endpoints, err = func() (b []string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]string, size)
		for i := 0; i < int(size); i++ {
			b[i], err = basic.ReadString(r)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Endpoints field: %s", err)
	}
	if s.SessionId, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read SessionId field: %s", err)
	}
	return s, nil
}
func WriteServiceInfo(s ServiceInfo, w io.Writer) (err error) {
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: %s", err)
	}
	if err := basic.WriteUint32(s.ServiceId, w); err != nil {
		return fmt.Errorf("failed to write ServiceId field: %s", err)
	}
	if err := basic.WriteString(s.MachineId, w); err != nil {
		return fmt.Errorf("failed to write MachineId field: %s", err)
	}
	if err := basic.WriteUint32(s.ProcessId, w); err != nil {
		return fmt.Errorf("failed to write ProcessId field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Endpoints)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.Endpoints {
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Endpoints field: %s", err)
	}
	if err := basic.WriteString(s.SessionId, w); err != nil {
		return fmt.Errorf("failed to write SessionId field: %s", err)
	}
	return nil
}

type LogMessage struct {
	Source     string
	Level      uint32
	Category   string
	Location   string
	Message    string
	Id         uint32
	Date       uint64
	SystemDate uint64
}

func ReadLogMessage(r io.Reader) (s LogMessage, err error) {
	if s.Source, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Source field: %s", err)
	}
	if s.Level, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read Level field: %s", err)
	}
	if s.Category, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Category field: %s", err)
	}
	if s.Location, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Location field: %s", err)
	}
	if s.Message, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Message field: %s", err)
	}
	if s.Id, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read Id field: %s", err)
	}
	if s.Date, err = basic.ReadUint64(r); err != nil {
		return s, fmt.Errorf("failed to read Date field: %s", err)
	}
	if s.SystemDate, err = basic.ReadUint64(r); err != nil {
		return s, fmt.Errorf("failed to read SystemDate field: %s", err)
	}
	return s, nil
}
func WriteLogMessage(s LogMessage, w io.Writer) (err error) {
	if err := basic.WriteString(s.Source, w); err != nil {
		return fmt.Errorf("failed to write Source field: %s", err)
	}
	if err := basic.WriteUint32(s.Level, w); err != nil {
		return fmt.Errorf("failed to write Level field: %s", err)
	}
	if err := basic.WriteString(s.Category, w); err != nil {
		return fmt.Errorf("failed to write Category field: %s", err)
	}
	if err := basic.WriteString(s.Location, w); err != nil {
		return fmt.Errorf("failed to write Location field: %s", err)
	}
	if err := basic.WriteString(s.Message, w); err != nil {
		return fmt.Errorf("failed to write Message field: %s", err)
	}
	if err := basic.WriteUint32(s.Id, w); err != nil {
		return fmt.Errorf("failed to write Id field: %s", err)
	}
	if err := basic.WriteUint64(s.Date, w); err != nil {
		return fmt.Errorf("failed to write Date field: %s", err)
	}
	if err := basic.WriteUint64(s.SystemDate, w); err != nil {
		return fmt.Errorf("failed to write SystemDate field: %s", err)
	}
	return nil
}
