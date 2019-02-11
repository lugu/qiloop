package stub

import (
	"bytes"
	"github.com/lugu/qiloop/meta/idl"
	"testing"
)

func helpGenerate(t *testing.T, input string) {
	pkg, err := idl.ParsePackage([]byte(input))
	if err != nil {
		panic(err)
	}
	if len(pkg.Types) == 0 {
		t.Fatalf("parse error: missing type")
	}
	var buf bytes.Buffer
	err = GeneratePackage(&buf, "", pkg)
	if err != nil {
		panic(err)
	}
	if string(buf.Bytes()) == "" {
		t.Fatalf("missing type declaration")
	}
}

func TestEmptyPackageName(t *testing.T) {
	var buf bytes.Buffer
	var pkg idl.PackageDeclaration
	err := GeneratePackage(&buf, "", &pkg)
	if err == nil {
		t.Errorf("shall not generate package without name")
	}
}

func TestGenerateEnum(t *testing.T) {
	input := `
	package test
	enum EnumType0
		valueA = 1
		valueB = 2
	end
	`
	helpGenerate(t, input)
}

func TestGenerateStruct(t *testing.T) {
	input := `
	package test
	struct B
		a: A
		b: str
	end
	struct A
		a: int32
		b: float32
	end
	`
	helpGenerate(t, input)
}

func TestGenerateEmptyInterface(t *testing.T) {
	input := `
	package test
	interface Empty
	end`
	helpGenerate(t, input)
}

func TestGenerateInterfaceMethod(t *testing.T) {
	input := `
	package test
	interface SimpleMethod
	    fn terminate(uid: uint32) //uid:6
	    fn property(P0: any) -> any //uid:5
	    fn authenticate(capability: Map<str,any>) -> Map<str,any> //uid:8
	    fn machineId() -> str //uid:108
	end`
	helpGenerate(t, input)
}

func TestGenerateInterfaceMethodWithStruct(t *testing.T) {
	input := `
	package test
	struct ServiceInfo
	    name: str
	    serviceId: uint32
	    machineId: str
	    processId: uint32
	    endpoints: Vec<str>
	    sessionId: str
	end
	interface StructMethod
	    fn service(P0: str) -> ServiceInfo //uid:100
	    fn updateServiceInfo(P0: ServiceInfo) //uid:105
	end`
	helpGenerate(t, input)
}

func TestGenerateInterfaceSignal(t *testing.T) {
	input := `
	package test
	interface SimpleSignal
	    sig serviceAdded(P0: uint32, P1: str) //uid:106
	    sig serviceRemoved(P0: uint32, P1: str) //uid:107
	end`
	helpGenerate(t, input)
}
