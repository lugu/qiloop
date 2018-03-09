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
func (p *Server) MetaObject(p0 uint32) (object.MetaObject, error) {
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

type PackageManager struct {
	session.Proxy
}

func NewPackageManager(ses session.Session, obj uint32) (*PackageManager, error) {
	proxy, err := ses.Proxy("PackageManager", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &PackageManager{proxy}, nil
}
func (p *PackageManager) RegisterEvent(p0 uint32, p1 uint32, p2 uint64) (uint64, error) {
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
func (p *PackageManager) UnregisterEvent(p0 uint32, p1 uint32, p2 uint64) error {
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
func (p *PackageManager) MetaObject(p0 uint32) (object.MetaObject, error) {
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
func (p *PackageManager) Terminate(p0 uint32) error {
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
func (p *PackageManager) Property(p0 value.Value) (value.Value, error) {
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
func (p *PackageManager) SetProperty(p0 value.Value, p1 value.Value) error {
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
func (p *PackageManager) Properties() ([]string, error) {
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
func (p *PackageManager) RegisterEventWithSignature(p0 uint32, p1 uint32, p2 uint64, p3 string) (uint64, error) {
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
func (p *PackageManager) IsStatsEnabled() (bool, error) {
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
func (p *PackageManager) EnableStats(p0 bool) error {
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
func (p *PackageManager) Stats() (map[uint32]MethodStatistics, error) {
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
func (p *PackageManager) ClearStats() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("clearStats", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearStats failed: %s", err)
	}
	return nil
}
func (p *PackageManager) IsTraceEnabled() (bool, error) {
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
func (p *PackageManager) EnableTrace(p0 bool) error {
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
func (p *PackageManager) Install(p0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("install", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call install failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse install response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) InstallCheckMd5(p0 string, p1 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteString(p1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p1: %s", err)
	}
	response, err := p.Call("installCheckMd5", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call installCheckMd5 failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse installCheckMd5 response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) _install(p0 string, p1 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteString(p1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p1: %s", err)
	}
	response, err := p.Call("_install", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _install failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _install response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) _install_0(p0 string, p1 string, p2 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteString(p1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p1: %s", err)
	}
	if err = basic.WriteString(p2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p2: %s", err)
	}
	response, err := p.Call("_install", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _install failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _install response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) _install_1(p0 string, p1 string, p2 string, p3 bool) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = basic.WriteString(p1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p1: %s", err)
	}
	if err = basic.WriteString(p2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p2: %s", err)
	}
	if err = basic.WriteBool(p3, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p3: %s", err)
	}
	response, err := p.Call("_install", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _install failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _install response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) HasPackage(p0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("hasPackage", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call hasPackage failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse hasPackage response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) Packages2() ([]PackageInfo2, error) {
	var err error
	var ret []PackageInfo2
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("packages2", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call packages2 failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (b []PackageInfo2, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]PackageInfo2, size)
		for i := 0; i < int(size); i++ {
			b[i], err = ReadPackageInfo2(buf)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse packages2 response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) Package2(p0 string) (PackageInfo2, error) {
	var err error
	var ret PackageInfo2
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("package2", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call package2 failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = ReadPackageInfo2(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse package2 response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) PackageIcon(p0 string) (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("packageIcon", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call packageIcon failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse packageIcon response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) RemovePkg(p0 string) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call("removePkg", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call removePkg failed: %s", err)
	}
	return nil
}
func (p *PackageManager) GetPackages() (value.Value, error) {
	var err error
	var ret value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("getPackages", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getPackages failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = value.NewValue(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse getPackages response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) Packages() ([]PackageInfo, error) {
	var err error
	var ret []PackageInfo
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("packages", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call packages failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (b []PackageInfo, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]PackageInfo, size)
		for i := 0; i < int(size); i++ {
			b[i], err = ReadPackageInfo(buf)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse packages response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) Package(p0 string) (PackageInfo, error) {
	var err error
	var ret PackageInfo
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("package", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call package failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = ReadPackageInfo(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse package response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) GetPackage(p0 string) (value.Value, error) {
	var err error
	var ret value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("getPackage", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getPackage failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = value.NewValue(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse getPackage response: %s", err)
	}
	return ret, nil
}
func (p *PackageManager) GetPackageIcon(p0 string) (value.Value, error) {
	var err error
	var ret value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call("getPackageIcon", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getPackageIcon failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = value.NewValue(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse getPackageIcon response: %s", err)
	}
	return ret, nil
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

type PackageInfo2 struct {
	Uuid         string
	Version      string
	Author       string
	Channel      string
	Organization string
	Date         string
	TypeVersion  string
	Installer    string
	Path         string
	Elems        map[string]value.Value
}

func ReadPackageInfo2(r io.Reader) (s PackageInfo2, err error) {
	if s.Uuid, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Uuid field: %s", err)
	}
	if s.Version, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Version field: %s", err)
	}
	if s.Author, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Author field: %s", err)
	}
	if s.Channel, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Channel field: %s", err)
	}
	if s.Organization, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Organization field: %s", err)
	}
	if s.Date, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Date field: %s", err)
	}
	if s.TypeVersion, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read TypeVersion field: %s", err)
	}
	if s.Installer, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Installer field: %s", err)
	}
	if s.Path, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Path field: %s", err)
	}
	if s.Elems, err = func() (m map[string]value.Value, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]value.Value, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := value.NewValue(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Elems field: %s", err)
	}
	return s, nil
}
func WritePackageInfo2(s PackageInfo2, w io.Writer) (err error) {
	if err := basic.WriteString(s.Uuid, w); err != nil {
		return fmt.Errorf("failed to write Uuid field: %s", err)
	}
	if err := basic.WriteString(s.Version, w); err != nil {
		return fmt.Errorf("failed to write Version field: %s", err)
	}
	if err := basic.WriteString(s.Author, w); err != nil {
		return fmt.Errorf("failed to write Author field: %s", err)
	}
	if err := basic.WriteString(s.Channel, w); err != nil {
		return fmt.Errorf("failed to write Channel field: %s", err)
	}
	if err := basic.WriteString(s.Organization, w); err != nil {
		return fmt.Errorf("failed to write Organization field: %s", err)
	}
	if err := basic.WriteString(s.Date, w); err != nil {
		return fmt.Errorf("failed to write Date field: %s", err)
	}
	if err := basic.WriteString(s.TypeVersion, w); err != nil {
		return fmt.Errorf("failed to write TypeVersion field: %s", err)
	}
	if err := basic.WriteString(s.Installer, w); err != nil {
		return fmt.Errorf("failed to write Installer field: %s", err)
	}
	if err := basic.WriteString(s.Path, w); err != nil {
		return fmt.Errorf("failed to write Path field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Elems)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.Elems {
			err = basic.WriteString(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = v.Write(w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Elems field: %s", err)
	}
	return nil
}

type BehaviorInfo struct {
	Path                   string
	Nature                 string
	LangToName             map[string]string
	LangToDesc             map[string]string
	Categories             string
	LangToTags             map[string][]string
	LangToTriggerSentences map[string][]string
	LangToLoadingResponses map[string][]string
	PurposeToCondition     map[string][]string
	Permissions            []string
	UserRequestable        bool
}

func ReadBehaviorInfo(r io.Reader) (s BehaviorInfo, err error) {
	if s.Path, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Path field: %s", err)
	}
	if s.Nature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Nature field: %s", err)
	}
	if s.LangToName, err = func() (m map[string]string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]string, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read LangToName field: %s", err)
	}
	if s.LangToDesc, err = func() (m map[string]string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]string, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read LangToDesc field: %s", err)
	}
	if s.Categories, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Categories field: %s", err)
	}
	if s.LangToTags, err = func() (m map[string][]string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string][]string, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := func() (b []string, err error) {
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
			}()
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read LangToTags field: %s", err)
	}
	if s.LangToTriggerSentences, err = func() (m map[string][]string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string][]string, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := func() (b []string, err error) {
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
			}()
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read LangToTriggerSentences field: %s", err)
	}
	if s.LangToLoadingResponses, err = func() (m map[string][]string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string][]string, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := func() (b []string, err error) {
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
			}()
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read LangToLoadingResponses field: %s", err)
	}
	if s.PurposeToCondition, err = func() (m map[string][]string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string][]string, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := func() (b []string, err error) {
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
			}()
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read PurposeToCondition field: %s", err)
	}
	if s.Permissions, err = func() (b []string, err error) {
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
		return s, fmt.Errorf("failed to read Permissions field: %s", err)
	}
	if s.UserRequestable, err = basic.ReadBool(r); err != nil {
		return s, fmt.Errorf("failed to read UserRequestable field: %s", err)
	}
	return s, nil
}
func WriteBehaviorInfo(s BehaviorInfo, w io.Writer) (err error) {
	if err := basic.WriteString(s.Path, w); err != nil {
		return fmt.Errorf("failed to write Path field: %s", err)
	}
	if err := basic.WriteString(s.Nature, w); err != nil {
		return fmt.Errorf("failed to write Nature field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.LangToName)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.LangToName {
			err = basic.WriteString(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write LangToName field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.LangToDesc)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.LangToDesc {
			err = basic.WriteString(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write LangToDesc field: %s", err)
	}
	if err := basic.WriteString(s.Categories, w); err != nil {
		return fmt.Errorf("failed to write Categories field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.LangToTags)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.LangToTags {
			err = basic.WriteString(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = func() error {
				err := basic.WriteUint32(uint32(len(v)), w)
				if err != nil {
					return fmt.Errorf("failed to write slice size: %s", err)
				}
				for _, v := range v {
					err = basic.WriteString(v, w)
					if err != nil {
						return fmt.Errorf("failed to write slice value: %s", err)
					}
				}
				return nil
			}()
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write LangToTags field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.LangToTriggerSentences)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.LangToTriggerSentences {
			err = basic.WriteString(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = func() error {
				err := basic.WriteUint32(uint32(len(v)), w)
				if err != nil {
					return fmt.Errorf("failed to write slice size: %s", err)
				}
				for _, v := range v {
					err = basic.WriteString(v, w)
					if err != nil {
						return fmt.Errorf("failed to write slice value: %s", err)
					}
				}
				return nil
			}()
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write LangToTriggerSentences field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.LangToLoadingResponses)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.LangToLoadingResponses {
			err = basic.WriteString(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = func() error {
				err := basic.WriteUint32(uint32(len(v)), w)
				if err != nil {
					return fmt.Errorf("failed to write slice size: %s", err)
				}
				for _, v := range v {
					err = basic.WriteString(v, w)
					if err != nil {
						return fmt.Errorf("failed to write slice value: %s", err)
					}
				}
				return nil
			}()
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write LangToLoadingResponses field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.PurposeToCondition)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.PurposeToCondition {
			err = basic.WriteString(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = func() error {
				err := basic.WriteUint32(uint32(len(v)), w)
				if err != nil {
					return fmt.Errorf("failed to write slice size: %s", err)
				}
				for _, v := range v {
					err = basic.WriteString(v, w)
					if err != nil {
						return fmt.Errorf("failed to write slice value: %s", err)
					}
				}
				return nil
			}()
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write PurposeToCondition field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Permissions)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.Permissions {
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Permissions field: %s", err)
	}
	if err := basic.WriteBool(s.UserRequestable, w); err != nil {
		return fmt.Errorf("failed to write UserRequestable field: %s", err)
	}
	return nil
}

type LanguageInfo struct {
	Path          string
	EngineName    string
	EngineVersion string
	Locale        string
	LangToName    map[string]string
}

func ReadLanguageInfo(r io.Reader) (s LanguageInfo, err error) {
	if s.Path, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Path field: %s", err)
	}
	if s.EngineName, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read EngineName field: %s", err)
	}
	if s.EngineVersion, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read EngineVersion field: %s", err)
	}
	if s.Locale, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Locale field: %s", err)
	}
	if s.LangToName, err = func() (m map[string]string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]string, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read LangToName field: %s", err)
	}
	return s, nil
}
func WriteLanguageInfo(s LanguageInfo, w io.Writer) (err error) {
	if err := basic.WriteString(s.Path, w); err != nil {
		return fmt.Errorf("failed to write Path field: %s", err)
	}
	if err := basic.WriteString(s.EngineName, w); err != nil {
		return fmt.Errorf("failed to write EngineName field: %s", err)
	}
	if err := basic.WriteString(s.EngineVersion, w); err != nil {
		return fmt.Errorf("failed to write EngineVersion field: %s", err)
	}
	if err := basic.WriteString(s.Locale, w); err != nil {
		return fmt.Errorf("failed to write Locale field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.LangToName)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.LangToName {
			err = basic.WriteString(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write LangToName field: %s", err)
	}
	return nil
}

type RobotRequirement struct {
	Model          string
	MinHeadVersion string
	MaxHeadVersion string
	MinBodyVersion string
	MaxBodyVersion string
}

func ReadRobotRequirement(r io.Reader) (s RobotRequirement, err error) {
	if s.Model, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Model field: %s", err)
	}
	if s.MinHeadVersion, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read MinHeadVersion field: %s", err)
	}
	if s.MaxHeadVersion, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read MaxHeadVersion field: %s", err)
	}
	if s.MinBodyVersion, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read MinBodyVersion field: %s", err)
	}
	if s.MaxBodyVersion, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read MaxBodyVersion field: %s", err)
	}
	return s, nil
}
func WriteRobotRequirement(s RobotRequirement, w io.Writer) (err error) {
	if err := basic.WriteString(s.Model, w); err != nil {
		return fmt.Errorf("failed to write Model field: %s", err)
	}
	if err := basic.WriteString(s.MinHeadVersion, w); err != nil {
		return fmt.Errorf("failed to write MinHeadVersion field: %s", err)
	}
	if err := basic.WriteString(s.MaxHeadVersion, w); err != nil {
		return fmt.Errorf("failed to write MaxHeadVersion field: %s", err)
	}
	if err := basic.WriteString(s.MinBodyVersion, w); err != nil {
		return fmt.Errorf("failed to write MinBodyVersion field: %s", err)
	}
	if err := basic.WriteString(s.MaxBodyVersion, w); err != nil {
		return fmt.Errorf("failed to write MaxBodyVersion field: %s", err)
	}
	return nil
}

type NaoqiRequirement struct {
	MinVersion string
	MaxVersion string
}

func ReadNaoqiRequirement(r io.Reader) (s NaoqiRequirement, err error) {
	if s.MinVersion, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read MinVersion field: %s", err)
	}
	if s.MaxVersion, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read MaxVersion field: %s", err)
	}
	return s, nil
}
func WriteNaoqiRequirement(s NaoqiRequirement, w io.Writer) (err error) {
	if err := basic.WriteString(s.MinVersion, w); err != nil {
		return fmt.Errorf("failed to write MinVersion field: %s", err)
	}
	if err := basic.WriteString(s.MaxVersion, w); err != nil {
		return fmt.Errorf("failed to write MaxVersion field: %s", err)
	}
	return nil
}

type PackageService struct {
	ExecStart string
	Name      string
	AutoRun   bool
}

func ReadPackageService(r io.Reader) (s PackageService, err error) {
	if s.ExecStart, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read ExecStart field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: %s", err)
	}
	if s.AutoRun, err = basic.ReadBool(r); err != nil {
		return s, fmt.Errorf("failed to read AutoRun field: %s", err)
	}
	return s, nil
}
func WritePackageService(s PackageService, w io.Writer) (err error) {
	if err := basic.WriteString(s.ExecStart, w); err != nil {
		return fmt.Errorf("failed to write ExecStart field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: %s", err)
	}
	if err := basic.WriteBool(s.AutoRun, w); err != nil {
		return fmt.Errorf("failed to write AutoRun field: %s", err)
	}
	return nil
}

type DialogInfo struct {
	TopicName   string
	TypeVersion string
	Topics      map[string]string
}

func ReadDialogInfo(r io.Reader) (s DialogInfo, err error) {
	if s.TopicName, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read TopicName field: %s", err)
	}
	if s.TypeVersion, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read TypeVersion field: %s", err)
	}
	if s.Topics, err = func() (m map[string]string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]string, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Topics field: %s", err)
	}
	return s, nil
}
func WriteDialogInfo(s DialogInfo, w io.Writer) (err error) {
	if err := basic.WriteString(s.TopicName, w); err != nil {
		return fmt.Errorf("failed to write TopicName field: %s", err)
	}
	if err := basic.WriteString(s.TypeVersion, w); err != nil {
		return fmt.Errorf("failed to write TypeVersion field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Topics)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.Topics {
			err = basic.WriteString(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Topics field: %s", err)
	}
	return nil
}

type PackageInfo struct {
	Uuid                 string
	Path                 string
	Version              string
	Channel              string
	Author               string
	Organization         string
	Date                 string
	TypeVersion          string
	LangToName           map[string]string
	LangToDesc           map[string]string
	SupportedLanguages   []string
	Behaviors            []BehaviorInfo
	Languages            []LanguageInfo
	Installer            string
	RobotRequirements    []RobotRequirement
	NaoqiRequirements    []NaoqiRequirement
	Services             []PackageService
	ExecutableFiles      []string
	Dialogs              []DialogInfo
	DescriptionLanguages []string
}

func ReadPackageInfo(r io.Reader) (s PackageInfo, err error) {
	if s.Uuid, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Uuid field: %s", err)
	}
	if s.Path, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Path field: %s", err)
	}
	if s.Version, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Version field: %s", err)
	}
	if s.Channel, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Channel field: %s", err)
	}
	if s.Author, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Author field: %s", err)
	}
	if s.Organization, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Organization field: %s", err)
	}
	if s.Date, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Date field: %s", err)
	}
	if s.TypeVersion, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read TypeVersion field: %s", err)
	}
	if s.LangToName, err = func() (m map[string]string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]string, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read LangToName field: %s", err)
	}
	if s.LangToDesc, err = func() (m map[string]string, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]string, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := basic.ReadString(r)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read LangToDesc field: %s", err)
	}
	if s.SupportedLanguages, err = func() (b []string, err error) {
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
		return s, fmt.Errorf("failed to read SupportedLanguages field: %s", err)
	}
	if s.Behaviors, err = func() (b []BehaviorInfo, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]BehaviorInfo, size)
		for i := 0; i < int(size); i++ {
			b[i], err = ReadBehaviorInfo(r)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Behaviors field: %s", err)
	}
	if s.Languages, err = func() (b []LanguageInfo, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]LanguageInfo, size)
		for i := 0; i < int(size); i++ {
			b[i], err = ReadLanguageInfo(r)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Languages field: %s", err)
	}
	if s.Installer, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Installer field: %s", err)
	}
	if s.RobotRequirements, err = func() (b []RobotRequirement, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]RobotRequirement, size)
		for i := 0; i < int(size); i++ {
			b[i], err = ReadRobotRequirement(r)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read RobotRequirements field: %s", err)
	}
	if s.NaoqiRequirements, err = func() (b []NaoqiRequirement, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]NaoqiRequirement, size)
		for i := 0; i < int(size); i++ {
			b[i], err = ReadNaoqiRequirement(r)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read NaoqiRequirements field: %s", err)
	}
	if s.Services, err = func() (b []PackageService, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]PackageService, size)
		for i := 0; i < int(size); i++ {
			b[i], err = ReadPackageService(r)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Services field: %s", err)
	}
	if s.ExecutableFiles, err = func() (b []string, err error) {
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
		return s, fmt.Errorf("failed to read ExecutableFiles field: %s", err)
	}
	if s.Dialogs, err = func() (b []DialogInfo, err error) {
		size, err := basic.ReadUint32(r)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]DialogInfo, size)
		for i := 0; i < int(size); i++ {
			b[i], err = ReadDialogInfo(r)
			if err != nil {
				return b, fmt.Errorf("failed to read slice value: %s", err)
			}
		}
		return b, nil
	}(); err != nil {
		return s, fmt.Errorf("failed to read Dialogs field: %s", err)
	}
	if s.DescriptionLanguages, err = func() (b []string, err error) {
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
		return s, fmt.Errorf("failed to read DescriptionLanguages field: %s", err)
	}
	return s, nil
}
func WritePackageInfo(s PackageInfo, w io.Writer) (err error) {
	if err := basic.WriteString(s.Uuid, w); err != nil {
		return fmt.Errorf("failed to write Uuid field: %s", err)
	}
	if err := basic.WriteString(s.Path, w); err != nil {
		return fmt.Errorf("failed to write Path field: %s", err)
	}
	if err := basic.WriteString(s.Version, w); err != nil {
		return fmt.Errorf("failed to write Version field: %s", err)
	}
	if err := basic.WriteString(s.Channel, w); err != nil {
		return fmt.Errorf("failed to write Channel field: %s", err)
	}
	if err := basic.WriteString(s.Author, w); err != nil {
		return fmt.Errorf("failed to write Author field: %s", err)
	}
	if err := basic.WriteString(s.Organization, w); err != nil {
		return fmt.Errorf("failed to write Organization field: %s", err)
	}
	if err := basic.WriteString(s.Date, w); err != nil {
		return fmt.Errorf("failed to write Date field: %s", err)
	}
	if err := basic.WriteString(s.TypeVersion, w); err != nil {
		return fmt.Errorf("failed to write TypeVersion field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.LangToName)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.LangToName {
			err = basic.WriteString(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write LangToName field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.LangToDesc)), w)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range s.LangToDesc {
			err = basic.WriteString(k, w)
			if err != nil {
				return fmt.Errorf("failed to write map key: %s", err)
			}
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write map value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write LangToDesc field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.SupportedLanguages)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.SupportedLanguages {
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write SupportedLanguages field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Behaviors)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.Behaviors {
			err = WriteBehaviorInfo(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Behaviors field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Languages)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.Languages {
			err = WriteLanguageInfo(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Languages field: %s", err)
	}
	if err := basic.WriteString(s.Installer, w); err != nil {
		return fmt.Errorf("failed to write Installer field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.RobotRequirements)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.RobotRequirements {
			err = WriteRobotRequirement(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write RobotRequirements field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.NaoqiRequirements)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.NaoqiRequirements {
			err = WriteNaoqiRequirement(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write NaoqiRequirements field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Services)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.Services {
			err = WritePackageService(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Services field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.ExecutableFiles)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.ExecutableFiles {
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write ExecutableFiles field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.Dialogs)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.Dialogs {
			err = WriteDialogInfo(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write Dialogs field: %s", err)
	}
	if err := func() error {
		err := basic.WriteUint32(uint32(len(s.DescriptionLanguages)), w)
		if err != nil {
			return fmt.Errorf("failed to write slice size: %s", err)
		}
		for _, v := range s.DescriptionLanguages {
			err = basic.WriteString(v, w)
			if err != nil {
				return fmt.Errorf("failed to write slice value: %s", err)
			}
		}
		return nil
	}(); err != nil {
		return fmt.Errorf("failed to write DescriptionLanguages field: %s", err)
	}
	return nil
}
