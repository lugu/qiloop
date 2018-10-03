package stub

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/type/object"
	"io"
)

func GenerateStubs(metas []object.MetaObject, packageName string, w io.Writer) error {
	for _, meta := range metas {
		err := GenerateStub(meta, packageName, meta.Description, w)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateStub(metaObj object.MetaObject, packageName, serviceName string, w io.Writer) error {
	file := jen.NewFile(packageName)
	file.PackageComment("file generated. DO NOT EDIT.")
	return fmt.Errorf("GenerateStub not yet implemented")
}
