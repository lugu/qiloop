// file generated. DO NOT EDIT.
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
