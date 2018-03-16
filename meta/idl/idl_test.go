package idl

import (
	"github.com/lugu/qiloop/type/object"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestServiceServer(t *testing.T) {
	var w strings.Builder
	if err := GenerateIDL(&w, "Server", object.MetaService0); err != nil {
		t.Errorf("failed to parse server: %s", err)
	}
	expected := `interface Server
	fn authenticate(P0: Map<str,Value>) -> Map<str,Value>
end
`
	if w.String() != expected {
		t.Errorf("Got:\n%s\nExpecting:\n%s\n", w.String(), expected)
	}
}

func TestObject(t *testing.T) {
	var w strings.Builder
	if err := GenerateIDL(&w, "Object", object.ObjectMetaObject); err != nil {
		t.Errorf("failed to parse server: %s", err)
	}
	expected := `interface Object
	fn registerEvent(P0: uint, P1: uint, P2: ulong) -> ulong
	fn unregisterEvent(P0: uint, P1: uint, P2: ulong)
	fn metaObject(P0: uint) -> MetaObject
	fn terminate(P0: uint)
	fn property(P0: Value) -> Value
	fn setProperty(P0: Value, P1: Value)
	fn properties() -> Vec<str>
	fn registerEventWithSignature(P0: uint, P1: uint, P2: ulong, P3: str) -> ulong
end
struct MetaMethodParameter
	name: str
	description: str
end
struct MetaMethod
	uid: uint
	returnSignature: str
	name: str
	parametersSignature: str
	description: str
	parameters: Vec<MetaMethodParameter>
	returnDescription: str
end
struct MetaSignal
	uid: uint
	name: str
	signature: str
end
struct MetaProperty
	uid: uint
	name: str
	signature: str
end
struct MetaObject
	methods: Map<uint,MetaMethod>
	signals: Map<uint,MetaSignal>
	properties: Map<uint,MetaProperty>
	description: str
end
`
	if w.String() != expected {
		t.Errorf("Got:\n%s\nExpecting:\n%s\n", w.String(), expected)
	}
}

func TestServiceDirectory(t *testing.T) {
	var w strings.Builder
	path := filepath.Join("testdata", "meta-object.bin")
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	metaObj, err := object.ReadMetaObject(file)
	if err := GenerateIDL(&w, "ServiceDirectory", metaObj); err != nil {
		t.Errorf("failed to parse server: %s", err)
	}
	expected := `interface ServiceDirectory
	fn registerEvent(P0: uint, P1: uint, P2: ulong) -> ulong
	fn unregisterEvent(P0: uint, P1: uint, P2: ulong)
	fn metaObject(P0: uint) -> MetaObject
	fn terminate(P0: uint)
	fn property(P0: Value) -> Value
	fn setProperty(P0: Value, P1: Value)
	fn properties() -> Vec<str>
	fn registerEventWithSignature(P0: uint, P1: uint, P2: ulong, P3: str) -> ulong
	fn isStatsEnabled() -> bool
	fn enableStats(P0: bool)
	fn stats() -> Map<uint,MethodStatistics>
	fn clearStats()
	fn isTraceEnabled() -> bool
	fn enableTrace(P0: bool)
	fn service(P0: str) -> ServiceInfo
	fn services() -> Vec<ServiceInfo>
	fn registerService(P0: ServiceInfo) -> uint
	fn unregisterService(P0: uint)
	fn serviceReady(P0: uint)
	fn updateServiceInfo(P0: ServiceInfo)
	fn machineId() -> str
	fn _socketOfService(P0: uint) -> interface
	sig traceObject(P0: EventTrace)
	sig serviceAdded(P0: uint, P1: str)
	sig serviceRemoved(P0: uint, P1: str)
end
struct MetaMethodParameter
	name: str
	description: str
end
struct MetaMethod
	uid: uint
	returnSignature: str
	name: str
	parametersSignature: str
	description: str
	parameters: Vec<MetaMethodParameter>
	returnDescription: str
end
struct MetaSignal
	uid: uint
	name: str
	signature: str
end
struct MetaProperty
	uid: uint
	name: str
	signature: str
end
struct MetaObject
	methods: Map<uint,MetaMethod>
	signals: Map<uint,MetaSignal>
	properties: Map<uint,MetaProperty>
	description: str
end
struct MinMaxSum
	minValue: float
	maxValue: float
	cumulatedValue: float
end
struct MethodStatistics
	count: uint
	wall: MinMaxSum
	user: MinMaxSum
	system: MinMaxSum
end
struct ServiceInfo
	name: str
	serviceId: uint
	machineId: str
	processId: uint
	endpoints: Vec<str>
	sessionId: str
end
struct timeval
	tv_sec: long
	tv_usec: long
end
struct EventTrace
	id: uint
	kind: int
	slotId: uint
	arguments: Value
	timestamp: timeval
	userUsTime: long
	systemUsTime: long
	callerContext: uint
	calleeContext: uint
end
`
	if w.String() != expected {
		t.Errorf("Got:\n%s\nExpecting:\n%s\n", w.String(), expected)
	}
}
