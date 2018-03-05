package services

import (
	bytes "bytes"
	fmt "fmt"
	basic "github.com/lugu/qiloop/basic"
	net "github.com/lugu/qiloop/net"
	value "github.com/lugu/qiloop/value"
	io "io"
)

type Server struct {
	net.Proxy
}

func NewServer(endpoint string, service uint32, obj uint32) (*Server, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &Server{proxy}, nil
	}
}
func (p *Server) MetaObject(p0 uint32) (MetaObject, error) {
	var err error
	var ret MetaObject
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call(2, buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call metaObject failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = ReadMetaObject(buf)
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
	response, err := p.Call(8, buf.Bytes())
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
	net.Proxy
}

func NewServiceDirectory(endpoint string, service uint32, obj uint32) (*ServiceDirectory, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ServiceDirectory{proxy}, nil
	}
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
	response, err := p.Call(0, buf.Bytes())
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
	_, err = p.Call(1, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterEvent failed: %s", err)
	}
	return nil
}
func (p *ServiceDirectory) MetaObject(p0 uint32) (MetaObject, error) {
	var err error
	var ret MetaObject
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call(2, buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call metaObject failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = ReadMetaObject(buf)
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
	_, err = p.Call(3, buf.Bytes())
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
	response, err := p.Call(5, buf.Bytes())
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
	_, err = p.Call(6, buf.Bytes())
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
	response, err := p.Call(7, buf.Bytes())
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
	response, err := p.Call(8, buf.Bytes())
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
	response, err := p.Call(80, buf.Bytes())
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
	_, err = p.Call(81, buf.Bytes())
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
	response, err := p.Call(82, buf.Bytes())
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
	_, err = p.Call(83, buf.Bytes())
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
	response, err := p.Call(84, buf.Bytes())
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
	_, err = p.Call(85, buf.Bytes())
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
	response, err := p.Call(100, buf.Bytes())
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
	response, err := p.Call(101, buf.Bytes())
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
	response, err := p.Call(102, buf.Bytes())
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
	_, err = p.Call(103, buf.Bytes())
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
	_, err = p.Call(104, buf.Bytes())
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
	_, err = p.Call(105, buf.Bytes())
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
	response, err := p.Call(108, buf.Bytes())
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
	net.Proxy
}

func NewLogManager(endpoint string, service uint32, obj uint32) (*LogManager, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &LogManager{proxy}, nil
	}
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
	response, err := p.Call(0, buf.Bytes())
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
	_, err = p.Call(1, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterEvent failed: %s", err)
	}
	return nil
}
func (p *LogManager) MetaObject(p0 uint32) (MetaObject, error) {
	var err error
	var ret MetaObject
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call(2, buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call metaObject failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = ReadMetaObject(buf)
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
	_, err = p.Call(3, buf.Bytes())
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
	response, err := p.Call(5, buf.Bytes())
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
	_, err = p.Call(6, buf.Bytes())
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
	response, err := p.Call(7, buf.Bytes())
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
	response, err := p.Call(8, buf.Bytes())
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
	response, err := p.Call(80, buf.Bytes())
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
	_, err = p.Call(81, buf.Bytes())
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
	response, err := p.Call(82, buf.Bytes())
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
	_, err = p.Call(83, buf.Bytes())
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
	response, err := p.Call(84, buf.Bytes())
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
	_, err = p.Call(85, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableTrace failed: %s", err)
	}
	return nil
}

type PackageManager struct {
	net.Proxy
}

func NewPackageManager(endpoint string, service uint32, obj uint32) (*PackageManager, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &PackageManager{proxy}, nil
	}
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
	response, err := p.Call(0, buf.Bytes())
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
	_, err = p.Call(1, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterEvent failed: %s", err)
	}
	return nil
}
func (p *PackageManager) MetaObject(p0 uint32) (MetaObject, error) {
	var err error
	var ret MetaObject
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call(2, buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call metaObject failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = ReadMetaObject(buf)
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
	_, err = p.Call(3, buf.Bytes())
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
	response, err := p.Call(5, buf.Bytes())
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
	_, err = p.Call(6, buf.Bytes())
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
	response, err := p.Call(7, buf.Bytes())
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
	response, err := p.Call(8, buf.Bytes())
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
	response, err := p.Call(80, buf.Bytes())
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
	_, err = p.Call(81, buf.Bytes())
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
	response, err := p.Call(82, buf.Bytes())
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
	_, err = p.Call(83, buf.Bytes())
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
	response, err := p.Call(84, buf.Bytes())
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
	_, err = p.Call(85, buf.Bytes())
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
	response, err := p.Call(100, buf.Bytes())
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
	response, err := p.Call(102, buf.Bytes())
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
	response, err := p.Call(103, buf.Bytes())
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
	response, err := p.Call(104, buf.Bytes())
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
	response, err := p.Call(105, buf.Bytes())
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
	response, err := p.Call(108, buf.Bytes())
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
	response, err := p.Call(109, buf.Bytes())
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
	response, err := p.Call(110, buf.Bytes())
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
	response, err := p.Call(111, buf.Bytes())
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
	_, err = p.Call(112, buf.Bytes())
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
	response, err := p.Call(115, buf.Bytes())
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
	response, err := p.Call(116, buf.Bytes())
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
	response, err := p.Call(117, buf.Bytes())
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
	response, err := p.Call(118, buf.Bytes())
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
	response, err := p.Call(119, buf.Bytes())
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

type ALServiceManager struct {
	net.Proxy
}

func NewALServiceManager(endpoint string, service uint32, obj uint32) (*ALServiceManager, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALServiceManager{proxy}, nil
	}
}
func (p *ALServiceManager) RegisterEvent(p0 uint32, p1 uint32, p2 uint64) (uint64, error) {
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
	response, err := p.Call(0, buf.Bytes())
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
func (p *ALServiceManager) UnregisterEvent(p0 uint32, p1 uint32, p2 uint64) error {
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
	_, err = p.Call(1, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterEvent failed: %s", err)
	}
	return nil
}
func (p *ALServiceManager) MetaObject(p0 uint32) (MetaObject, error) {
	var err error
	var ret MetaObject
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call(2, buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call metaObject failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = ReadMetaObject(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse metaObject response: %s", err)
	}
	return ret, nil
}
func (p *ALServiceManager) Terminate(p0 uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call(3, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call terminate failed: %s", err)
	}
	return nil
}
func (p *ALServiceManager) Property(p0 value.Value) (value.Value, error) {
	var err error
	var ret value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = p0.Write(buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call(5, buf.Bytes())
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
func (p *ALServiceManager) SetProperty(p0 value.Value, p1 value.Value) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = p0.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	if err = p1.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize p1: %s", err)
	}
	_, err = p.Call(6, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setProperty failed: %s", err)
	}
	return nil
}
func (p *ALServiceManager) Properties() ([]string, error) {
	var err error
	var ret []string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call(7, buf.Bytes())
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
func (p *ALServiceManager) RegisterEventWithSignature(p0 uint32, p1 uint32, p2 uint64, p3 string) (uint64, error) {
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
	response, err := p.Call(8, buf.Bytes())
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
func (p *ALServiceManager) IsStatsEnabled() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call(80, buf.Bytes())
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
func (p *ALServiceManager) EnableStats(p0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call(81, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableStats failed: %s", err)
	}
	return nil
}
func (p *ALServiceManager) Stats() (map[uint32]MethodStatistics, error) {
	var err error
	var ret map[uint32]MethodStatistics
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call(82, buf.Bytes())
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
func (p *ALServiceManager) ClearStats() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call(83, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearStats failed: %s", err)
	}
	return nil
}
func (p *ALServiceManager) IsTraceEnabled() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call(84, buf.Bytes())
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
func (p *ALServiceManager) EnableTrace(p0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call(85, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableTrace failed: %s", err)
	}
	return nil
}
func (p *ALServiceManager) Start(p0 string) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call(100, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call start failed: %s", err)
	}
	return nil
}
func (p *ALServiceManager) Restart(p0 string) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call(101, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call restart failed: %s", err)
	}
	return nil
}
func (p *ALServiceManager) Stop(p0 string) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return fmt.Errorf("failed to serialize p0: %s", err)
	}
	_, err = p.Call(102, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call stop failed: %s", err)
	}
	return nil
}
func (p *ALServiceManager) StopAllServices() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call(103, buf.Bytes())
	if err != nil {
		return fmt.Errorf("call stopAllServices failed: %s", err)
	}
	return nil
}
func (p *ALServiceManager) IsServiceRunning(p0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call(104, buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call isServiceRunning failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse isServiceRunning response: %s", err)
	}
	return ret, nil
}
func (p *ALServiceManager) ServiceMemoryUsage(p0 string) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call(105, buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call serviceMemoryUsage failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadUint64(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse serviceMemoryUsage response: %s", err)
	}
	return ret, nil
}
func (p *ALServiceManager) Services() ([]ServiceProcessInfo, error) {
	var err error
	var ret []ServiceProcessInfo
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call(106, buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call services failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = func() (b []ServiceProcessInfo, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return b, fmt.Errorf("failed to read slice size: %s", err)
		}
		b = make([]ServiceProcessInfo, size)
		for i := 0; i < int(size); i++ {
			b[i], err = ReadServiceProcessInfo(buf)
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
func (p *ALServiceManager) Service(p0 string) (ServiceProcessInfo, error) {
	var err error
	var ret ServiceProcessInfo
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call(107, buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call service failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = ReadServiceProcessInfo(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse service response: %s", err)
	}
	return ret, nil
}
func (p *ALServiceManager) StartService(p0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call(110, buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call startService failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse startService response: %s", err)
	}
	return ret, nil
}
func (p *ALServiceManager) RestartService(p0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call(111, buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call restartService failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse restartService response: %s", err)
	}
	return ret, nil
}
func (p *ALServiceManager) StopService(p0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(p0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize p0: %s", err)
	}
	response, err := p.Call(112, buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call stopService failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse stopService response: %s", err)
	}
	return ret, nil
}

type ALFileManager struct {
	net.Proxy
}

func NewALFileManager(endpoint string, service uint32, obj uint32) (*ALFileManager, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALFileManager{proxy}, nil
	}
}

type ALMemory struct {
	net.Proxy
}

func NewALMemory(endpoint string, service uint32, obj uint32) (*ALMemory, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALMemory{proxy}, nil
	}
}

type ALLogger struct {
	net.Proxy
}

func NewALLogger(endpoint string, service uint32, obj uint32) (*ALLogger, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALLogger{proxy}, nil
	}
}

type ALPreferences struct {
	net.Proxy
}

func NewALPreferences(endpoint string, service uint32, obj uint32) (*ALPreferences, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALPreferences{proxy}, nil
	}
}

type ALLauncher struct {
	net.Proxy
}

func NewALLauncher(endpoint string, service uint32, obj uint32) (*ALLauncher, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALLauncher{proxy}, nil
	}
}

type ALDebug struct {
	net.Proxy
}

func NewALDebug(endpoint string, service uint32, obj uint32) (*ALDebug, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALDebug{proxy}, nil
	}
}

type ALPreferenceManager struct {
	net.Proxy
}

func NewALPreferenceManager(endpoint string, service uint32, obj uint32) (*ALPreferenceManager, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALPreferenceManager{proxy}, nil
	}
}

type ALNotificationManager struct {
	net.Proxy
}

func NewALNotificationManager(endpoint string, service uint32, obj uint32) (*ALNotificationManager, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALNotificationManager{proxy}, nil
	}
}

type _ALNotificationAdder struct {
	net.Proxy
}

func New_ALNotificationAdder(endpoint string, service uint32, obj uint32) (*_ALNotificationAdder, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_ALNotificationAdder{proxy}, nil
	}
}

type ALResourceManager struct {
	net.Proxy
}

func NewALResourceManager(endpoint string, service uint32, obj uint32) (*ALResourceManager, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALResourceManager{proxy}, nil
	}
}

type ALRobotModel struct {
	net.Proxy
}

func NewALRobotModel(endpoint string, service uint32, obj uint32) (*ALRobotModel, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALRobotModel{proxy}, nil
	}
}

type ALSonar struct {
	net.Proxy
}

func NewALSonar(endpoint string, service uint32, obj uint32) (*ALSonar, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALSonar{proxy}, nil
	}
}

type ALFsr struct {
	net.Proxy
}

func NewALFsr(endpoint string, service uint32, obj uint32) (*ALFsr, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALFsr{proxy}, nil
	}
}

type ALSensors struct {
	net.Proxy
}

func NewALSensors(endpoint string, service uint32, obj uint32) (*ALSensors, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALSensors{proxy}, nil
	}
}

type ALBodyTemperature struct {
	net.Proxy
}

func NewALBodyTemperature(endpoint string, service uint32, obj uint32) (*ALBodyTemperature, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALBodyTemperature{proxy}, nil
	}
}

type ALMotion struct {
	net.Proxy
}

func NewALMotion(endpoint string, service uint32, obj uint32) (*ALMotion, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALMotion{proxy}, nil
	}
}

type ALTouch struct {
	net.Proxy
}

func NewALTouch(endpoint string, service uint32, obj uint32) (*ALTouch, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALTouch{proxy}, nil
	}
}

type ALRobotPosture struct {
	net.Proxy
}

func NewALRobotPosture(endpoint string, service uint32, obj uint32) (*ALRobotPosture, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALRobotPosture{proxy}, nil
	}
}

type ALMotionRecorder struct {
	net.Proxy
}

func NewALMotionRecorder(endpoint string, service uint32, obj uint32) (*ALMotionRecorder, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALMotionRecorder{proxy}, nil
	}
}

type ALLeds struct {
	net.Proxy
}

func NewALLeds(endpoint string, service uint32, obj uint32) (*ALLeds, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALLeds{proxy}, nil
	}
}

type ALWorldRepresentation struct {
	net.Proxy
}

func NewALWorldRepresentation(endpoint string, service uint32, obj uint32) (*ALWorldRepresentation, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALWorldRepresentation{proxy}, nil
	}
}

type ALKnowledgeManager struct {
	net.Proxy
}

func NewALKnowledgeManager(endpoint string, service uint32, obj uint32) (*ALKnowledgeManager, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALKnowledgeManager{proxy}, nil
	}
}

type ALKnowledge struct {
	net.Proxy
}

func NewALKnowledge(endpoint string, service uint32, obj uint32) (*ALKnowledge, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALKnowledge{proxy}, nil
	}
}

type ALAudioPlayer struct {
	net.Proxy
}

func NewALAudioPlayer(endpoint string, service uint32, obj uint32) (*ALAudioPlayer, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALAudioPlayer{proxy}, nil
	}
}

type ALTextToSpeech struct {
	net.Proxy
}

func NewALTextToSpeech(endpoint string, service uint32, obj uint32) (*ALTextToSpeech, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALTextToSpeech{proxy}, nil
	}
}

type ALBattery struct {
	net.Proxy
}

func NewALBattery(endpoint string, service uint32, obj uint32) (*ALBattery, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALBattery{proxy}, nil
	}
}

type ALFrameManager struct {
	net.Proxy
}

func NewALFrameManager(endpoint string, service uint32, obj uint32) (*ALFrameManager, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALFrameManager{proxy}, nil
	}
}

type ALPythonBridge struct {
	net.Proxy
}

func NewALPythonBridge(endpoint string, service uint32, obj uint32) (*ALPythonBridge, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALPythonBridge{proxy}, nil
	}
}

type ALVideoDevice struct {
	net.Proxy
}

func NewALVideoDevice(endpoint string, service uint32, obj uint32) (*ALVideoDevice, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALVideoDevice{proxy}, nil
	}
}

type ALRedBallDetection struct {
	net.Proxy
}

func NewALRedBallDetection(endpoint string, service uint32, obj uint32) (*ALRedBallDetection, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALRedBallDetection{proxy}, nil
	}
}

type ALVisionRecognition struct {
	net.Proxy
}

func NewALVisionRecognition(endpoint string, service uint32, obj uint32) (*ALVisionRecognition, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALVisionRecognition{proxy}, nil
	}
}

type ALBehaviorManager struct {
	net.Proxy
}

func NewALBehaviorManager(endpoint string, service uint32, obj uint32) (*ALBehaviorManager, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALBehaviorManager{proxy}, nil
	}
}

type _ALMovementScheduler struct {
	net.Proxy
}

func New_ALMovementScheduler(endpoint string, service uint32, obj uint32) (*_ALMovementScheduler, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_ALMovementScheduler{proxy}, nil
	}
}

type ALAnimationPlayer struct {
	net.Proxy
}

func NewALAnimationPlayer(endpoint string, service uint32, obj uint32) (*ALAnimationPlayer, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALAnimationPlayer{proxy}, nil
	}
}

type ALSpeakingMovement struct {
	net.Proxy
}

func NewALSpeakingMovement(endpoint string, service uint32, obj uint32) (*ALSpeakingMovement, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALSpeakingMovement{proxy}, nil
	}
}

type ALAnimatedSpeech struct {
	net.Proxy
}

func NewALAnimatedSpeech(endpoint string, service uint32, obj uint32) (*ALAnimatedSpeech, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALAnimatedSpeech{proxy}, nil
	}
}

type ALColorBlobDetection struct {
	net.Proxy
}

func NewALColorBlobDetection(endpoint string, service uint32, obj uint32) (*ALColorBlobDetection, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALColorBlobDetection{proxy}, nil
	}
}

type ALVisualSpaceHistory struct {
	net.Proxy
}

func NewALVisualSpaceHistory(endpoint string, service uint32, obj uint32) (*ALVisualSpaceHistory, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALVisualSpaceHistory{proxy}, nil
	}
}

type _RedBallConverter struct {
	net.Proxy
}

func New_RedBallConverter(endpoint string, service uint32, obj uint32) (*_RedBallConverter, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_RedBallConverter{proxy}, nil
	}
}

type _FaceConverter struct {
	net.Proxy
}

func New_FaceConverter(endpoint string, service uint32, obj uint32) (*_FaceConverter, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_FaceConverter{proxy}, nil
	}
}

type _NaoMarkConverter1 struct {
	net.Proxy
}

func New_NaoMarkConverter1(endpoint string, service uint32, obj uint32) (*_NaoMarkConverter1, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_NaoMarkConverter1{proxy}, nil
	}
}

type _NaoMarkConverter2 struct {
	net.Proxy
}

func New_NaoMarkConverter2(endpoint string, service uint32, obj uint32) (*_NaoMarkConverter2, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_NaoMarkConverter2{proxy}, nil
	}
}

type _LogoConverter struct {
	net.Proxy
}

func New_LogoConverter(endpoint string, service uint32, obj uint32) (*_LogoConverter, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_LogoConverter{proxy}, nil
	}
}

type _SoundConverter struct {
	net.Proxy
}

func New_SoundConverter(endpoint string, service uint32, obj uint32) (*_SoundConverter, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_SoundConverter{proxy}, nil
	}
}

type _PeopleConverter struct {
	net.Proxy
}

func New_PeopleConverter(endpoint string, service uint32, obj uint32) (*_PeopleConverter, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_PeopleConverter{proxy}, nil
	}
}

type _Smoother struct {
	net.Proxy
}

func New_Smoother(endpoint string, service uint32, obj uint32) (*_Smoother, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_Smoother{proxy}, nil
	}
}

type _MotionMove struct {
	net.Proxy
}

func New_MotionMove(endpoint string, service uint32, obj uint32) (*_MotionMove, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_MotionMove{proxy}, nil
	}
}

type _NavigationMove struct {
	net.Proxy
}

func New_NavigationMove(endpoint string, service uint32, obj uint32) (*_NavigationMove, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_NavigationMove{proxy}, nil
	}
}

type _RoundSearcher struct {
	net.Proxy
}

func New_RoundSearcher(endpoint string, service uint32, obj uint32) (*_RoundSearcher, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_RoundSearcher{proxy}, nil
	}
}

type _HeadLooker struct {
	net.Proxy
}

func New_HeadLooker(endpoint string, service uint32, obj uint32) (*_HeadLooker, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_HeadLooker{proxy}, nil
	}
}

type _WholeBodyLooker struct {
	net.Proxy
}

func New_WholeBodyLooker(endpoint string, service uint32, obj uint32) (*_WholeBodyLooker, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_WholeBodyLooker{proxy}, nil
	}
}

type _ArmsLooker struct {
	net.Proxy
}

func New_ArmsLooker(endpoint string, service uint32, obj uint32) (*_ArmsLooker, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_ArmsLooker{proxy}, nil
	}
}

type _ALTargetManager struct {
	net.Proxy
}

func New_ALTargetManager(endpoint string, service uint32, obj uint32) (*_ALTargetManager, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_ALTargetManager{proxy}, nil
	}
}

type ALTracker struct {
	net.Proxy
}

func NewALTracker(endpoint string, service uint32, obj uint32) (*ALTracker, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALTracker{proxy}, nil
	}
}

type ALModularity struct {
	net.Proxy
}

func NewALModularity(endpoint string, service uint32, obj uint32) (*ALModularity, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALModularity{proxy}, nil
	}
}

type ALNavigation struct {
	net.Proxy
}

func NewALNavigation(endpoint string, service uint32, obj uint32) (*ALNavigation, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALNavigation{proxy}, nil
	}
}

type ALMovementDetection struct {
	net.Proxy
}

func NewALMovementDetection(endpoint string, service uint32, obj uint32) (*ALMovementDetection, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALMovementDetection{proxy}, nil
	}
}

type ALSegmentation3D struct {
	net.Proxy
}

func NewALSegmentation3D(endpoint string, service uint32, obj uint32) (*ALSegmentation3D, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALSegmentation3D{proxy}, nil
	}
}

type _UserIdentification struct {
	net.Proxy
}

func New_UserIdentification(endpoint string, service uint32, obj uint32) (*_UserIdentification, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_UserIdentification{proxy}, nil
	}
}

type ALPeoplePerception struct {
	net.Proxy
}

func NewALPeoplePerception(endpoint string, service uint32, obj uint32) (*ALPeoplePerception, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALPeoplePerception{proxy}, nil
	}
}

type ALSittingPeopleDetection struct {
	net.Proxy
}

func NewALSittingPeopleDetection(endpoint string, service uint32, obj uint32) (*ALSittingPeopleDetection, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALSittingPeopleDetection{proxy}, nil
	}
}

type ALEngagementZones struct {
	net.Proxy
}

func NewALEngagementZones(endpoint string, service uint32, obj uint32) (*ALEngagementZones, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALEngagementZones{proxy}, nil
	}
}

type ALGazeAnalysis struct {
	net.Proxy
}

func NewALGazeAnalysis(endpoint string, service uint32, obj uint32) (*ALGazeAnalysis, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALGazeAnalysis{proxy}, nil
	}
}

type ALWavingDetection struct {
	net.Proxy
}

func NewALWavingDetection(endpoint string, service uint32, obj uint32) (*ALWavingDetection, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALWavingDetection{proxy}, nil
	}
}

type ALCloseObjectDetection struct {
	net.Proxy
}

func NewALCloseObjectDetection(endpoint string, service uint32, obj uint32) (*ALCloseObjectDetection, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALCloseObjectDetection{proxy}, nil
	}
}

type _ALFastPersonTracking struct {
	net.Proxy
}

func New_ALFastPersonTracking(endpoint string, service uint32, obj uint32) (*_ALFastPersonTracking, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_ALFastPersonTracking{proxy}, nil
	}
}

type _ALFindPersonHead struct {
	net.Proxy
}

func New_ALFindPersonHead(endpoint string, service uint32, obj uint32) (*_ALFindPersonHead, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_ALFindPersonHead{proxy}, nil
	}
}

type ALVisualCompass struct {
	net.Proxy
}

func NewALVisualCompass(endpoint string, service uint32, obj uint32) (*ALVisualCompass, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALVisualCompass{proxy}, nil
	}
}

type ALLocalization struct {
	net.Proxy
}

func NewALLocalization(endpoint string, service uint32, obj uint32) (*ALLocalization, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALLocalization{proxy}, nil
	}
}

type ALPanoramaCompass struct {
	net.Proxy
}

func NewALPanoramaCompass(endpoint string, service uint32, obj uint32) (*ALPanoramaCompass, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALPanoramaCompass{proxy}, nil
	}
}

type ALUserInfo struct {
	net.Proxy
}

func NewALUserInfo(endpoint string, service uint32, obj uint32) (*ALUserInfo, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALUserInfo{proxy}, nil
	}
}

type ALUserSession struct {
	net.Proxy
}

func NewALUserSession(endpoint string, service uint32, obj uint32) (*ALUserSession, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALUserSession{proxy}, nil
	}
}

type _ALUserSessionCompat struct {
	net.Proxy
}

func New_ALUserSessionCompat(endpoint string, service uint32, obj uint32) (*_ALUserSessionCompat, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_ALUserSessionCompat{proxy}, nil
	}
}

type ALThinkingExpression struct {
	net.Proxy
}

func NewALThinkingExpression(endpoint string, service uint32, obj uint32) (*ALThinkingExpression, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALThinkingExpression{proxy}, nil
	}
}

type ALBasicAwareness struct {
	net.Proxy
}

func NewALBasicAwareness(endpoint string, service uint32, obj uint32) (*ALBasicAwareness, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALBasicAwareness{proxy}, nil
	}
}

type ALBackgroundMovement struct {
	net.Proxy
}

func NewALBackgroundMovement(endpoint string, service uint32, obj uint32) (*ALBackgroundMovement, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALBackgroundMovement{proxy}, nil
	}
}

type ALListeningMovement struct {
	net.Proxy
}

func NewALListeningMovement(endpoint string, service uint32, obj uint32) (*ALListeningMovement, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALListeningMovement{proxy}, nil
	}
}

type _ConditionChecker_maxime_2388_0 struct {
	net.Proxy
}

func New_ConditionChecker_maxime_2388_0(endpoint string, service uint32, obj uint32) (*_ConditionChecker_maxime_2388_0, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_ConditionChecker_maxime_2388_0{proxy}, nil
	}
}

type ALExpressionWatcher struct {
	net.Proxy
}

func NewALExpressionWatcher(endpoint string, service uint32, obj uint32) (*ALExpressionWatcher, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALExpressionWatcher{proxy}, nil
	}
}

type _LifeReporter struct {
	net.Proxy
}

func New_LifeReporter(endpoint string, service uint32, obj uint32) (*_LifeReporter, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_LifeReporter{proxy}, nil
	}
}

type _ConditionChecker_maxime_2388_1 struct {
	net.Proxy
}

func New_ConditionChecker_maxime_2388_1(endpoint string, service uint32, obj uint32) (*_ConditionChecker_maxime_2388_1, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_ConditionChecker_maxime_2388_1{proxy}, nil
	}
}

type _LaunchpadPluginActivities struct {
	net.Proxy
}

func New_LaunchpadPluginActivities(endpoint string, service uint32, obj uint32) (*_LaunchpadPluginActivities, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_LaunchpadPluginActivities{proxy}, nil
	}
}

type _LaunchpadPluginTiming struct {
	net.Proxy
}

func New_LaunchpadPluginTiming(endpoint string, service uint32, obj uint32) (*_LaunchpadPluginTiming, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_LaunchpadPluginTiming{proxy}, nil
	}
}

type _LaunchpadPluginTemperature struct {
	net.Proxy
}

func New_LaunchpadPluginTemperature(endpoint string, service uint32, obj uint32) (*_LaunchpadPluginTemperature, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_LaunchpadPluginTemperature{proxy}, nil
	}
}

type _LaunchpadPluginHardwareInfo struct {
	net.Proxy
}

func New_LaunchpadPluginHardwareInfo(endpoint string, service uint32, obj uint32) (*_LaunchpadPluginHardwareInfo, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_LaunchpadPluginHardwareInfo{proxy}, nil
	}
}

type _LaunchpadPluginBattery struct {
	net.Proxy
}

func New_LaunchpadPluginBattery(endpoint string, service uint32, obj uint32) (*_LaunchpadPluginBattery, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_LaunchpadPluginBattery{proxy}, nil
	}
}

type _LaunchpadPluginPosture struct {
	net.Proxy
}

func New_LaunchpadPluginPosture(endpoint string, service uint32, obj uint32) (*_LaunchpadPluginPosture, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_LaunchpadPluginPosture{proxy}, nil
	}
}

type _ConditionChecker_maxime_2388_2 struct {
	net.Proxy
}

func New_ConditionChecker_maxime_2388_2(endpoint string, service uint32, obj uint32) (*_ConditionChecker_maxime_2388_2, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_ConditionChecker_maxime_2388_2{proxy}, nil
	}
}

type _ExitInteractive struct {
	net.Proxy
}

func New_ExitInteractive(endpoint string, service uint32, obj uint32) (*_ExitInteractive, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_ExitInteractive{proxy}, nil
	}
}

type _UserSessionManager struct {
	net.Proxy
}

func New_UserSessionManager(endpoint string, service uint32, obj uint32) (*_UserSessionManager, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &_UserSessionManager{proxy}, nil
	}
}

type ALAutonomousLife struct {
	net.Proxy
}

func NewALAutonomousLife(endpoint string, service uint32, obj uint32) (*ALAutonomousLife, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALAutonomousLife{proxy}, nil
	}
}

type ALDialog struct {
	net.Proxy
}

func NewALDialog(endpoint string, service uint32, obj uint32) (*ALDialog, error) {
	if conn, err := net.NewClient(endpoint); err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", endpoint, err)
	} else {
		proxy := net.NewProxy(conn, service, obj)
		return &ALDialog{proxy}, nil
	}
}

type MetaMethodParameter struct {
	Name        string
	Description string
}

func ReadMetaMethodParameter(r io.Reader) (s MetaMethodParameter, err error) {
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: %s", err)
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Description field: %s", err)
	}
	return s, nil
}
func WriteMetaMethodParameter(s MetaMethodParameter, w io.Writer) (err error) {
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: %s", err)
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("failed to write Description field: %s", err)
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
		return s, fmt.Errorf("failed to read Uid field: %s", err)
	}
	if s.ReturnSignature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read ReturnSignature field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: %s", err)
	}
	if s.ParametersSignature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read ParametersSignature field: %s", err)
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Description field: %s", err)
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
		return s, fmt.Errorf("failed to read Parameters field: %s", err)
	}
	if s.ReturnDescription, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read ReturnDescription field: %s", err)
	}
	return s, nil
}
func WriteMetaMethod(s MetaMethod, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("failed to write Uid field: %s", err)
	}
	if err := basic.WriteString(s.ReturnSignature, w); err != nil {
		return fmt.Errorf("failed to write ReturnSignature field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: %s", err)
	}
	if err := basic.WriteString(s.ParametersSignature, w); err != nil {
		return fmt.Errorf("failed to write ParametersSignature field: %s", err)
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("failed to write Description field: %s", err)
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
		return fmt.Errorf("failed to write Parameters field: %s", err)
	}
	if err := basic.WriteString(s.ReturnDescription, w); err != nil {
		return fmt.Errorf("failed to write ReturnDescription field: %s", err)
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
		return s, fmt.Errorf("failed to read Uid field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: %s", err)
	}
	if s.Signature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Signature field: %s", err)
	}
	return s, nil
}
func WriteMetaSignal(s MetaSignal, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("failed to write Uid field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: %s", err)
	}
	if err := basic.WriteString(s.Signature, w); err != nil {
		return fmt.Errorf("failed to write Signature field: %s", err)
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
		return s, fmt.Errorf("failed to read Uid field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: %s", err)
	}
	if s.Signature, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Signature field: %s", err)
	}
	return s, nil
}
func WriteMetaProperty(s MetaProperty, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Uid, w); err != nil {
		return fmt.Errorf("failed to write Uid field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: %s", err)
	}
	if err := basic.WriteString(s.Signature, w); err != nil {
		return fmt.Errorf("failed to write Signature field: %s", err)
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
		return s, fmt.Errorf("failed to read Methods field: %s", err)
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
		return s, fmt.Errorf("failed to read Signals field: %s", err)
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
		return s, fmt.Errorf("failed to read Properties field: %s", err)
	}
	if s.Description, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Description field: %s", err)
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
		return fmt.Errorf("failed to write Methods field: %s", err)
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
		return fmt.Errorf("failed to write Signals field: %s", err)
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
		return fmt.Errorf("failed to write Properties field: %s", err)
	}
	if err := basic.WriteString(s.Description, w); err != nil {
		return fmt.Errorf("failed to write Description field: %s", err)
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

type ServiceProcessInfo struct {
	Running   bool
	Name      string
	ExecStart string
	Autorun   bool
}

func ReadServiceProcessInfo(r io.Reader) (s ServiceProcessInfo, err error) {
	if s.Running, err = basic.ReadBool(r); err != nil {
		return s, fmt.Errorf("failed to read Running field: %s", err)
	}
	if s.Name, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read Name field: %s", err)
	}
	if s.ExecStart, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("failed to read ExecStart field: %s", err)
	}
	if s.Autorun, err = basic.ReadBool(r); err != nil {
		return s, fmt.Errorf("failed to read Autorun field: %s", err)
	}
	return s, nil
}
func WriteServiceProcessInfo(s ServiceProcessInfo, w io.Writer) (err error) {
	if err := basic.WriteBool(s.Running, w); err != nil {
		return fmt.Errorf("failed to write Running field: %s", err)
	}
	if err := basic.WriteString(s.Name, w); err != nil {
		return fmt.Errorf("failed to write Name field: %s", err)
	}
	if err := basic.WriteString(s.ExecStart, w); err != nil {
		return fmt.Errorf("failed to write ExecStart field: %s", err)
	}
	if err := basic.WriteBool(s.Autorun, w); err != nil {
		return fmt.Errorf("failed to write Autorun field: %s", err)
	}
	return nil
}
