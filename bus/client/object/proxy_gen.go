// file generated. DO NOT EDIT.
package object

import (
	"bytes"
	"fmt"
	bus "github.com/lugu/qiloop/bus"
	basic "github.com/lugu/qiloop/type/basic"
	object "github.com/lugu/qiloop/type/object"
	value "github.com/lugu/qiloop/type/value"
	"io"
	"log"
)

type NewServices struct {
	session bus.Session
}

func Services(s bus.Session) NewServices {
	return NewServices{session: s}
}

type Object interface {
	object.Object
	bus.Proxy
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
func (s NewServices) Object() (Object, error) {
	return NewObject(s.session, 1)
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
func (p *ObjectProxy) IsStatsEnabled() (bool, error) {
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
func (p *ObjectProxy) EnableStats(P0 bool) error {
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
func (p *ObjectProxy) Stats() (map[uint32]MethodStatistics, error) {
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
func (p *ObjectProxy) ClearStats() error {
	var err error
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(make([]byte, 0))
	_, err = p.Call("clearStats", buf.Bytes())
	if err != nil {
		return fmt.Errorf("call clearStats failed: %s", err)
	}
	return nil
}
func (p *ObjectProxy) IsTraceEnabled() (bool, error) {
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
func (p *ObjectProxy) EnableTrace(P0 bool) error {
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
func (p *ObjectProxy) SignalTraceObject(cancel chan int) (chan struct {
	P0 EventTrace
}, error) {
	signalID, err := p.SignalUid("traceObject")
	if err != nil {
		return nil, fmt.Errorf("signal %s not available: %s", "traceObject", err)
	}

	handlerID := uint64(signalID)<<32 + 1 // FIXME: read it from proxy
	_, err = p.RegisterEvent(p.ObjectID(), signalID, handlerID)
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
				// connection lost or cancellation.
				close(ch)
				p.UnregisterEvent(p.ObjectID(), signalID, handlerID)
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
				log.Printf("failed to unmarshall tuple: %s", err)
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
		return s, fmt.Errorf("failed to read MinValue field: " + err.Error())
	}
	if s.MaxValue, err = basic.ReadFloat32(r); err != nil {
		return s, fmt.Errorf("failed to read MaxValue field: " + err.Error())
	}
	if s.CumulatedValue, err = basic.ReadFloat32(r); err != nil {
		return s, fmt.Errorf("failed to read CumulatedValue field: " + err.Error())
	}
	return s, nil
}
func WriteMinMaxSum(s MinMaxSum, w io.Writer) (err error) {
	if err := basic.WriteFloat32(s.MinValue, w); err != nil {
		return fmt.Errorf("failed to write MinValue field: " + err.Error())
	}
	if err := basic.WriteFloat32(s.MaxValue, w); err != nil {
		return fmt.Errorf("failed to write MaxValue field: " + err.Error())
	}
	if err := basic.WriteFloat32(s.CumulatedValue, w); err != nil {
		return fmt.Errorf("failed to write CumulatedValue field: " + err.Error())
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
		return s, fmt.Errorf("failed to read Count field: " + err.Error())
	}
	if s.Wall, err = ReadMinMaxSum(r); err != nil {
		return s, fmt.Errorf("failed to read Wall field: " + err.Error())
	}
	if s.User, err = ReadMinMaxSum(r); err != nil {
		return s, fmt.Errorf("failed to read User field: " + err.Error())
	}
	if s.System, err = ReadMinMaxSum(r); err != nil {
		return s, fmt.Errorf("failed to read System field: " + err.Error())
	}
	return s, nil
}
func WriteMethodStatistics(s MethodStatistics, w io.Writer) (err error) {
	if err := basic.WriteUint32(s.Count, w); err != nil {
		return fmt.Errorf("failed to write Count field: " + err.Error())
	}
	if err := WriteMinMaxSum(s.Wall, w); err != nil {
		return fmt.Errorf("failed to write Wall field: " + err.Error())
	}
	if err := WriteMinMaxSum(s.User, w); err != nil {
		return fmt.Errorf("failed to write User field: " + err.Error())
	}
	if err := WriteMinMaxSum(s.System, w); err != nil {
		return fmt.Errorf("failed to write System field: " + err.Error())
	}
	return nil
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