// file generated. DO NOT EDIT.
package services

import (
	bytes "bytes"
	fmt "fmt"
	bus "github.com/lugu/qiloop/bus"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
	io "io"
)

type ServerProxy struct {
	bus.Proxy
}

func NewServer(ses bus.Session, obj uint32) (Server, error) {
	proxy, err := ses.Proxy("Server", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &ServerProxy{proxy}, nil
}
func (p *ServerProxy) Authenticate(P0 map[string]value.Value) (map[string]value.Value, error) {
	var err error
	var ret map[string]value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = func() error {
		err := basic.WriteUint32(uint32(len(P0)), buf)
		if err != nil {
			return fmt.Errorf("failed to write map size: %s", err)
		}
		for k, v := range P0 {
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
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
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

type ObjectProxy struct {
	bus.Proxy
}

func NewObject(ses bus.Session, obj uint32) (Object, error) {
	proxy, err := ses.Proxy("Object", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &ObjectProxy{proxy}, nil
}
func (p *ObjectProxy) RegisterEvent(P0 uint32, P1 uint32, P2 uint64) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteUint32(P1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteUint64(P2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P2: %s", err)
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
func (p *ObjectProxy) UnregisterEvent(P0 uint32, P1 uint32, P2 uint64) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteUint32(P1, buf); err != nil {
		return fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteUint64(P2, buf); err != nil {
		return fmt.Errorf("failed to serialize P2: %s", err)
	}
	_, err = p.Call("unregisterEvent", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterEvent failed: %s", err)
	}
	return nil
}
func (p *ObjectProxy) MetaObject(P0 uint32) (object.MetaObject, error) {
	var err error
	var ret object.MetaObject
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
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
func (p *ObjectProxy) Terminate(P0 uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("terminate", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call terminate failed: %s", err)
	}
	return nil
}
func (p *ObjectProxy) Property(P0 value.Value) (value.Value, error) {
	var err error
	var ret value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = P0.Write(buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
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
func (p *ObjectProxy) SetProperty(P0 value.Value, P1 value.Value) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = P0.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = P1.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize P1: %s", err)
	}
	_, err = p.Call("setProperty", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setProperty failed: %s", err)
	}
	return nil
}
func (p *ObjectProxy) Properties() ([]string, error) {
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
func (p *ObjectProxy) RegisterEventWithSignature(P0 uint32, P1 uint32, P2 uint64, P3 string) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteUint32(P1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteUint64(P2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P2: %s", err)
	}
	if err = basic.WriteString(P3, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P3: %s", err)
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

type ALTabletServiceProxy struct {
	bus.Proxy
}

func NewALTabletService(ses bus.Session, obj uint32) (ALTabletService, error) {
	proxy, err := ses.Proxy("ALTabletService", obj)
	if err != nil {
		return nil, fmt.Errorf("failed to contact service: %s", err)
	}
	return &ALTabletServiceProxy{proxy}, nil
}
func (p *ALTabletServiceProxy) RegisterEvent(P0 uint32, P1 uint32, P2 uint64) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteUint32(P1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteUint64(P2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P2: %s", err)
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
func (p *ALTabletServiceProxy) UnregisterEvent(P0 uint32, P1 uint32, P2 uint64) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteUint32(P1, buf); err != nil {
		return fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteUint64(P2, buf); err != nil {
		return fmt.Errorf("failed to serialize P2: %s", err)
	}
	_, err = p.Call("unregisterEvent", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call unregisterEvent failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) MetaObject(P0 uint32) (object.MetaObject, error) {
	var err error
	var ret object.MetaObject
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
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
func (p *ALTabletServiceProxy) Terminate(P0 uint32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("terminate", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call terminate failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) Property(P0 value.Value) (value.Value, error) {
	var err error
	var ret value.Value
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = P0.Write(buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
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
func (p *ALTabletServiceProxy) SetProperty(P0 value.Value, P1 value.Value) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = P0.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = P1.Write(buf); err != nil {
		return fmt.Errorf("failed to serialize P1: %s", err)
	}
	_, err = p.Call("setProperty", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setProperty failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) Properties() ([]string, error) {
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
func (p *ALTabletServiceProxy) RegisterEventWithSignature(P0 uint32, P1 uint32, P2 uint64, P3 string) (uint64, error) {
	var err error
	var ret uint64
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteUint32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteUint32(P1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteUint64(P2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P2: %s", err)
	}
	if err = basic.WriteString(P3, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P3: %s", err)
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
func (p *ALTabletServiceProxy) IsStatsEnabled() (bool, error) {
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
func (p *ALTabletServiceProxy) EnableStats(P0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("enableStats", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableStats failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) Stats() (map[uint32]MethodStatistics, error) {
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
func (p *ALTabletServiceProxy) ClearStats() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("clearStats", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearStats failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) IsTraceEnabled() (bool, error) {
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
func (p *ALTabletServiceProxy) EnableTrace(P0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("enableTrace", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableTrace failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) ShowWebview() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("showWebview", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call showWebview failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse showWebview response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) ShowWebview_0(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("showWebview", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call showWebview failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse showWebview response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) HideWebview() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("hideWebview", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call hideWebview failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse hideWebview response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) CleanWebview() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("cleanWebview", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call cleanWebview failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _clearWebviewCache(P0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("_clearWebviewCache", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _clearWebviewCache failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) LoadUrl(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("loadUrl", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call loadUrl failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse loadUrl response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) ReloadPage(P0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("reloadPage", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call reloadPage failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) LoadApplication(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("loadApplication", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call loadApplication failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse loadApplication response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) ExecuteJS(P0 string) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("executeJS", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call executeJS failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) GetOnTouchScaleFactor() (float32, error) {
	var err error
	var ret float32
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("getOnTouchScaleFactor", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getOnTouchScaleFactor failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadFloat32(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse getOnTouchScaleFactor response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) SetOnTouchWebviewScaleFactor(P0 float32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteFloat32(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("setOnTouchWebviewScaleFactor", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call setOnTouchWebviewScaleFactor failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _setAnimatedCrossWalkView(P0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("_setAnimatedCrossWalkView", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _setAnimatedCrossWalkView failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _setDebugCrossWalkViewEnable(P0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("_setDebugCrossWalkViewEnable", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _setDebugCrossWalkViewEnable failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) PlayVideo(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("playVideo", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call playVideo failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse playVideo response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) ResumeVideo() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("resumeVideo", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call resumeVideo failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse resumeVideo response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) PauseVideo() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("pauseVideo", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call pauseVideo failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse pauseVideo response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) StopVideo() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("stopVideo", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call stopVideo failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse stopVideo response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) GetVideoPosition() (int32, error) {
	var err error
	var ret int32
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("getVideoPosition", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getVideoPosition failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadInt32(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse getVideoPosition response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) GetVideoLength() (int32, error) {
	var err error
	var ret int32
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("getVideoLength", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getVideoLength failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadInt32(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse getVideoLength response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) PreLoadImage(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("preLoadImage", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call preLoadImage failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse preLoadImage response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) ShowImage(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("showImage", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call showImage failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse showImage response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) ShowImageNoCache(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("showImageNoCache", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call showImageNoCache failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse showImageNoCache response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) HideImage() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("hideImage", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call hideImage failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) ResumeGif() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("resumeGif", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call resumeGif failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) PauseGif() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("pauseGif", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call pauseGif failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) SetBackgroundColor(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("setBackgroundColor", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call setBackgroundColor failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse setBackgroundColor response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _startAnimation(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("_startAnimation", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _startAnimation failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _startAnimation response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _stopAnimation() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("_stopAnimation", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _stopAnimation failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) Hide() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("hide", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call hide failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) SetBrightness(P0 float32) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteFloat32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("setBrightness", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call setBrightness failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse setBrightness response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) GetBrightness() (float32, error) {
	var err error
	var ret float32
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("getBrightness", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getBrightness failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadFloat32(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse getBrightness response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) TurnScreenOn(P0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("turnScreenOn", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call turnScreenOn failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) GoToSleep() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("goToSleep", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call goToSleep failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) WakeUp() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("wakeUp", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call wakeUp failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _displayToast(P0 string, P1 int32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteInt32(P1, buf); err != nil {
		return fmt.Errorf("failed to serialize P1: %s", err)
	}
	_, err = p.Call("_displayToast", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _displayToast failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _displayToast_0(P0 string, P1 int32, P2 int32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteInt32(P1, buf); err != nil {
		return fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteInt32(P2, buf); err != nil {
		return fmt.Errorf("failed to serialize P2: %s", err)
	}
	_, err = p.Call("_displayToast", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _displayToast failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) GetWifiStatus() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("getWifiStatus", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getWifiStatus failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse getWifiStatus response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) EnableWifi() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("enableWifi", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableWifi failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) DisableWifi() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("disableWifi", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call disableWifi failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) ForgetWifi(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("forgetWifi", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call forgetWifi failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse forgetWifi response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) ConnectWifi(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("connectWifi", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call connectWifi failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse connectWifi response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) DisconnectWifi() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("disconnectWifi", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call disconnectWifi failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse disconnectWifi response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) ConfigureWifi(P0 string, P1 string, P2 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteString(P1, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteString(P2, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P2: %s", err)
	}
	response, err := p.Call("configureWifi", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call configureWifi failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse configureWifi response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) GetWifiMacAddress() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("getWifiMacAddress", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getWifiMacAddress failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse getWifiMacAddress response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) ShowInputDialog(P0 string, P1 string, P2 string, P3 string) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteString(P1, buf); err != nil {
		return fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteString(P2, buf); err != nil {
		return fmt.Errorf("failed to serialize P2: %s", err)
	}
	if err = basic.WriteString(P3, buf); err != nil {
		return fmt.Errorf("failed to serialize P3: %s", err)
	}
	_, err = p.Call("showInputDialog", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call showInputDialog failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) ShowInputTextDialog(P0 string, P1 string, P2 string) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteString(P1, buf); err != nil {
		return fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteString(P2, buf); err != nil {
		return fmt.Errorf("failed to serialize P2: %s", err)
	}
	_, err = p.Call("showInputTextDialog", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call showInputTextDialog failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) ShowInputTextDialog_0(P0 string, P1 string, P2 string, P3 string, P4 int32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteString(P1, buf); err != nil {
		return fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteString(P2, buf); err != nil {
		return fmt.Errorf("failed to serialize P2: %s", err)
	}
	if err = basic.WriteString(P3, buf); err != nil {
		return fmt.Errorf("failed to serialize P3: %s", err)
	}
	if err = basic.WriteInt32(P4, buf); err != nil {
		return fmt.Errorf("failed to serialize P4: %s", err)
	}
	_, err = p.Call("showInputTextDialog", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call showInputTextDialog failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) ShowInputDialog_0(P0 string, P1 string, P2 string, P3 string, P4 string, P5 int32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	if err = basic.WriteString(P1, buf); err != nil {
		return fmt.Errorf("failed to serialize P1: %s", err)
	}
	if err = basic.WriteString(P2, buf); err != nil {
		return fmt.Errorf("failed to serialize P2: %s", err)
	}
	if err = basic.WriteString(P3, buf); err != nil {
		return fmt.Errorf("failed to serialize P3: %s", err)
	}
	if err = basic.WriteString(P4, buf); err != nil {
		return fmt.Errorf("failed to serialize P4: %s", err)
	}
	if err = basic.WriteInt32(P5, buf); err != nil {
		return fmt.Errorf("failed to serialize P5: %s", err)
	}
	_, err = p.Call("showInputDialog", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call showInputDialog failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) HideDialog() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("hideDialog", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call hideDialog failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) SetKeyboard(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("setKeyboard", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call setKeyboard failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse setKeyboard response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) GetAvailableKeyboards() ([]string, error) {
	var err error
	var ret []string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("getAvailableKeyboards", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call getAvailableKeyboards failed: %s", err)
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
		return ret, fmt.Errorf("failed to parse getAvailableKeyboards response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) SetTabletLanguage(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("setTabletLanguage", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call setTabletLanguage failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse setTabletLanguage response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _openSettings() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("_openSettings", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _openSettings failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) SetVolume(P0 int32) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteInt32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("setVolume", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call setVolume failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse setVolume response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _setDebugEnabled(P0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("_setDebugEnabled", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _setDebugEnabled failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _setTimeZone(P0 string) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("_setTimeZone", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _setTimeZone failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _setStackTraceDepth(P0 int32) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteInt32(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("_setStackTraceDepth", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _setStackTraceDepth failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _setStackTraceDepth response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _ping() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("_ping", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _ping failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _ping response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) RobotIp() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("robotIp", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call robotIp failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse robotIp response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) Version() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("version", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call version failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse version response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _firmwareVersion() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("_firmwareVersion", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _firmwareVersion failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _firmwareVersion response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) ResetTablet() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("resetTablet", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call resetTablet failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _enableResetTablet(P0 bool) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteBool(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("_enableResetTablet", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _enableResetTablet failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _cancelReset() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("_cancelReset", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _cancelReset failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _setOpenGLState(P0 int32) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteInt32(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("_setOpenGLState", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _setOpenGLState failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _updateFirmware(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("_updateFirmware", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _updateFirmware failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _updateFirmware response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _getTabletSerialno() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("_getTabletSerialno", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _getTabletSerialno failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _getTabletSerialno response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _getTabletModelName() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("_getTabletModelName", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _getTabletModelName failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _getTabletModelName response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _uninstallApps() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("_uninstallApps", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _uninstallApps failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _uninstallLauncher() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("_uninstallLauncher", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _uninstallLauncher failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _uninstallBrowser() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("_uninstallBrowser", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _uninstallBrowser failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _wipeData() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("_wipeData", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _wipeData failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _fingerPrint() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("_fingerPrint", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _fingerPrint failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _fingerPrint response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _powerOff() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("_powerOff", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _powerOff failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _installApk(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("_installApk", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _installApk failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _installApk response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _installSystemApk(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("_installSystemApk", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _installSystemApk failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _installSystemApk response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _launchApk(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("_launchApk", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _launchApk failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _launchApk response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _removeApk(P0 string) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("_removeApk", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _removeApk failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _listApks() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("_listApks", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _listApks failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _listApks response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _stopApk(P0 string) error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return fmt.Errorf("failed to serialize P0: %s", err)
	}
	_, err = p.Call("_stopApk", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _stopApk failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) _isApkExist(P0 string) (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("_isApkExist", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _isApkExist failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _isApkExist response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _getApkVersion(P0 string) (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("_getApkVersion", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _getApkVersion failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _getApkVersion response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _getApkVersionCode(P0 string) (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	if err = basic.WriteString(P0, buf); err != nil {
		return ret, fmt.Errorf("failed to serialize P0: %s", err)
	}
	response, err := p.Call("_getApkVersionCode", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _getApkVersionCode failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _getApkVersionCode response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _purgeInstallTabletUpdater() (bool, error) {
	var err error
	var ret bool
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("_purgeInstallTabletUpdater", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _purgeInstallTabletUpdater failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadBool(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _purgeInstallTabletUpdater response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _test() (string, error) {
	var err error
	var ret string
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	response, err := p.Call("_test", buf.Bytes())
	if err != nil {
		return ret, fmt.Errorf("call _test failed: %s", err)
	}
	buf = bytes.NewBuffer(response)
	ret, err = basic.ReadString(buf)
	if err != nil {
		return ret, fmt.Errorf("failed to parse _test response: %s", err)
	}
	return ret, nil
}
func (p *ALTabletServiceProxy) _crash() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("_crash", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call _crash failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) EnableWebviewTouch() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("enableWebviewTouch", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call enableWebviewTouch failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) DisableWebviewTouch() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("disableWebviewTouch", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call disableWebviewTouch failed: %s", err)
	}
	return nil
}
func (p *ALTabletServiceProxy) SignalTraceObject(cancel chan int) (chan struct {
	P0 EventTrace
}, error) {
	signalID, err := p.SignalUid("traceObject")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "traceObject", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "traceObject", err)
	}
	ch := make(chan struct {
		P0 EventTrace
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "traceObject", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 EventTrace
			}, err error) {
				s.P0, err = ReadEventTrace(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnTouch(cancel chan int) (chan struct {
	P0 float32
	P1 float32
}, error) {
	signalID, err := p.SignalUid("onTouch")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onTouch", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onTouch", err)
	}
	ch := make(chan struct {
		P0 float32
		P1 float32
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onTouch", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 float32
				P1 float32
			}, err error) {
				s.P0, err = basic.ReadFloat32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				s.P1, err = basic.ReadFloat32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnTouchDown(cancel chan int) (chan struct {
	P0 float32
	P1 float32
}, error) {
	signalID, err := p.SignalUid("onTouchDown")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onTouchDown", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onTouchDown", err)
	}
	ch := make(chan struct {
		P0 float32
		P1 float32
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onTouchDown", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 float32
				P1 float32
			}, err error) {
				s.P0, err = basic.ReadFloat32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				s.P1, err = basic.ReadFloat32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnTouchDownRatio(cancel chan int) (chan struct {
	P0 float32
	P1 float32
	P2 string
}, error) {
	signalID, err := p.SignalUid("onTouchDownRatio")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onTouchDownRatio", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onTouchDownRatio", err)
	}
	ch := make(chan struct {
		P0 float32
		P1 float32
		P2 string
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onTouchDownRatio", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 float32
				P1 float32
				P2 string
			}, err error) {
				s.P0, err = basic.ReadFloat32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				s.P1, err = basic.ReadFloat32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				s.P2, err = basic.ReadString(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnTouchUp(cancel chan int) (chan struct {
	P0 float32
	P1 float32
}, error) {
	signalID, err := p.SignalUid("onTouchUp")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onTouchUp", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onTouchUp", err)
	}
	ch := make(chan struct {
		P0 float32
		P1 float32
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onTouchUp", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 float32
				P1 float32
			}, err error) {
				s.P0, err = basic.ReadFloat32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				s.P1, err = basic.ReadFloat32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnTouchMove(cancel chan int) (chan struct {
	P0 float32
	P1 float32
}, error) {
	signalID, err := p.SignalUid("onTouchMove")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onTouchMove", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onTouchMove", err)
	}
	ch := make(chan struct {
		P0 float32
		P1 float32
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onTouchMove", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 float32
				P1 float32
			}, err error) {
				s.P0, err = basic.ReadFloat32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				s.P1, err = basic.ReadFloat32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalVideoFinished(cancel chan int) (chan struct{}, error) {
	signalID, err := p.SignalUid("videoFinished")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "videoFinished", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "videoFinished", err)
	}
	ch := make(chan struct{})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "videoFinished", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct{}, err error) {
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalVideoStarted(cancel chan int) (chan struct{}, error) {
	signalID, err := p.SignalUid("videoStarted")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "videoStarted", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "videoStarted", err)
	}
	ch := make(chan struct{})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "videoStarted", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct{}, err error) {
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnPageStarted(cancel chan int) (chan struct{}, error) {
	signalID, err := p.SignalUid("onPageStarted")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onPageStarted", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onPageStarted", err)
	}
	ch := make(chan struct{})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onPageStarted", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct{}, err error) {
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnPageFinished(cancel chan int) (chan struct{}, error) {
	signalID, err := p.SignalUid("onPageFinished")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onPageFinished", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onPageFinished", err)
	}
	ch := make(chan struct{})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onPageFinished", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct{}, err error) {
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnLoadPageError(cancel chan int) (chan struct {
	P0 int32
	P1 string
	P2 string
}, error) {
	signalID, err := p.SignalUid("onLoadPageError")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onLoadPageError", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onLoadPageError", err)
	}
	ch := make(chan struct {
		P0 int32
		P1 string
		P2 string
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onLoadPageError", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 int32
				P1 string
				P2 string
			}, err error) {
				s.P0, err = basic.ReadInt32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				s.P1, err = basic.ReadString(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				s.P2, err = basic.ReadString(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnWifiStatusChange(cancel chan int) (chan struct {
	P0 string
}, error) {
	signalID, err := p.SignalUid("onWifiStatusChange")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onWifiStatusChange", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onWifiStatusChange", err)
	}
	ch := make(chan struct {
		P0 string
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onWifiStatusChange", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 string
			}, err error) {
				s.P0, err = basic.ReadString(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnConsoleMessage(cancel chan int) (chan struct {
	P0 string
}, error) {
	signalID, err := p.SignalUid("onConsoleMessage")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onConsoleMessage", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onConsoleMessage", err)
	}
	ch := make(chan struct {
		P0 string
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onConsoleMessage", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 string
			}, err error) {
				s.P0, err = basic.ReadString(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnInputText(cancel chan int) (chan struct {
	P0 int32
	P1 string
}, error) {
	signalID, err := p.SignalUid("onInputText")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onInputText", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onInputText", err)
	}
	ch := make(chan struct {
		P0 int32
		P1 string
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onInputText", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 int32
				P1 string
			}, err error) {
				s.P0, err = basic.ReadInt32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				s.P1, err = basic.ReadString(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnApkInstalled(cancel chan int) (chan struct {
	P0 string
}, error) {
	signalID, err := p.SignalUid("onApkInstalled")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onApkInstalled", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onApkInstalled", err)
	}
	ch := make(chan struct {
		P0 string
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onApkInstalled", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 string
			}, err error) {
				s.P0, err = basic.ReadString(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnSystemApkInstalled(cancel chan int) (chan struct {
	P0 string
}, error) {
	signalID, err := p.SignalUid("onSystemApkInstalled")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onSystemApkInstalled", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onSystemApkInstalled", err)
	}
	ch := make(chan struct {
		P0 string
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onSystemApkInstalled", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 string
			}, err error) {
				s.P0, err = basic.ReadString(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnSystemApkInstallError(cancel chan int) (chan struct {
	P0 int32
}, error) {
	signalID, err := p.SignalUid("onSystemApkInstallError")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onSystemApkInstallError", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onSystemApkInstallError", err)
	}
	ch := make(chan struct {
		P0 int32
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onSystemApkInstallError", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 int32
			}, err error) {
				s.P0, err = basic.ReadInt32(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnImageLoaded(cancel chan int) (chan struct{}, error) {
	signalID, err := p.SignalUid("onImageLoaded")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onImageLoaded", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onImageLoaded", err)
	}
	ch := make(chan struct{})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onImageLoaded", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct{}, err error) {
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
}
func (p *ALTabletServiceProxy) SignalOnJSEvent(cancel chan int) (chan struct {
	P0 string
}, error) {
	signalID, err := p.SignalUid("onJSEvent")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "onJSEvent", err)
	}

	id, err := p.RegisterEvent(p.ObjectID(), signalID, uint64(signalID)<<32+1)
	if err != nil {
		return nil, fmt.Errorf("failed to register event for %s: %s", "onJSEvent", err)
	}
	ch := make(chan struct {
		P0 string
	})
	chPay, err := p.SubscribeID(signalID, cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to request signal: %s", err)
	}
	go func() {
		for {
			payload, ok := <-chPay
			if !ok {
				close(ch) // upstream is closed.
				err = p.UnregisterEvent(p.ObjectID(), signalID, id)
				if err != nil {
					// FIXME: implement proper logging.
					fmt.Printf("failed to unregister event %s: %s", "onJSEvent", err)
				}
				return
			}
			buf := bytes.NewBuffer(payload)
			_ = buf // discard unused variable error
			e, err := func() (s struct {
				P0 string
			}, err error) {
				s.P0, err = basic.ReadString(buf)
				if err != nil {
					return s, fmt.Errorf("failed to read tuple member: %s", err)
				}
				return s, nil
			}()
			if err != nil {
				fmt.Errorf("failed to unmarshall tuple: %s", err)
				continue
			}
			ch <- e
		}
	}()
	return ch, nil
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

type timeval struct {
	Tv_sec  int64
	Tv_usec int64
}

func Readtimeval(r io.Reader) (s timeval, err error) {
	if s.Tv_sec, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read Tv_sec field: %s", err)
	}
	if s.Tv_usec, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read Tv_usec field: %s", err)
	}
	return s, nil
}
func Writetimeval(s timeval, w io.Writer) (err error) {
	if err := basic.WriteInt64(s.Tv_sec, w); err != nil {
		return fmt.Errorf("failed to write Tv_sec field: %s", err)
	}
	if err := basic.WriteInt64(s.Tv_usec, w); err != nil {
		return fmt.Errorf("failed to write Tv_usec field: %s", err)
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
		return s, fmt.Errorf("failed to read Id field: %s", err)
	}
	if s.Kind, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("failed to read Kind field: %s", err)
	}
	if s.SlotId, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read SlotId field: %s", err)
	}
	if s.Arguments, err = value.NewValue(r); err != nil {
		return s, fmt.Errorf("failed to read Arguments field: %s", err)
	}
	if s.Timestamp, err = Readtimeval(r); err != nil {
		return s, fmt.Errorf("failed to read Timestamp field: %s", err)
	}
	if s.UserUsTime, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read UserUsTime field: %s", err)
	}
	if s.SystemUsTime, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("failed to read SystemUsTime field: %s", err)
	}
	if s.CallerContext, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read CallerContext field: %s", err)
	}
	if s.CalleeContext, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("failed to read CalleeContext field: %s", err)
	}
	return s, nil
}
func WriteEventTrace(s EventTrace, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Id, w); err != nil {
		return fmt.Errorf("failed to write Id field: %s", err)
	}
	if err := basic.WriteInt32(s.Kind, w); err != nil {
		return fmt.Errorf("failed to write Kind field: %s", err)
	}
	if err := basic.WriteUint32(s.SlotId, w); err != nil {
		return fmt.Errorf("failed to write SlotId field: %s", err)
	}
	if err := s.Arguments.Write(w); err != nil {
		return fmt.Errorf("failed to write Arguments field: %s", err)
	}
	if err := Writetimeval(s.Timestamp, w); err != nil {
		return fmt.Errorf("failed to write Timestamp field: %s", err)
	}
	if err := basic.WriteInt64(s.UserUsTime, w); err != nil {
		return fmt.Errorf("failed to write UserUsTime field: %s", err)
	}
	if err := basic.WriteInt64(s.SystemUsTime, w); err != nil {
		return fmt.Errorf("failed to write SystemUsTime field: %s", err)
	}
	if err := basic.WriteUint32(s.CallerContext, w); err != nil {
		return fmt.Errorf("failed to write CallerContext field: %s", err)
	}
	if err := basic.WriteUint32(s.CalleeContext, w); err != nil {
		return fmt.Errorf("failed to write CalleeContext field: %s", err)
	}
	return nil
}
