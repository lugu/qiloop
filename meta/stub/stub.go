package stub

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/type/object"
	"io"
)

func GenerateStub(metaObj object.MetaObject, packageName, serviceName string, w io.Writer) error {
	file := jen.NewFile(packageName)
	file.PackageComment("file generated. DO NOT EDIT.")
	return fmt.Errorf("GenerateStub not yet implemented")
}
