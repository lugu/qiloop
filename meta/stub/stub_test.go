package stub

import (
	"bytes"
	"github.com/lugu/qiloop/meta/idl"
	"testing"
)

func helpGenerate(t *testing.T, input string) string {
	pkg, err := idl.ParsePackage([]byte(input))
	if err != nil {
		panic(err)
	}
	if len(pkg.Types) == 0 {
		t.Fatalf("parse error: missing type")
	}
	var buf bytes.Buffer
	err = GeneratePackage(&buf, pkg)
	if err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}

func TestEmptyPackageName(t *testing.T) {
	var buf bytes.Buffer
	var pkg idl.PackageDeclaration
	err := GeneratePackage(&buf, &pkg)
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
	output := helpGenerate(t, input)
	if output == "" {
		t.Fatalf("missing type declaration")
	}
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
	output := helpGenerate(t, input)
	if output == "" {
		t.Fatalf("missing type declaration")
	}
}

func TestGenerateEmptyInterface(t *testing.T) {
	input := `
	package test
	interface Empty
	end`
	output := helpGenerate(t, input)
	if output == "" {
		t.Fatalf("missing type declaration")
	}
}
