package stub

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/signature"
	"io"
)

func GeneratePackage(pkg *idl.PackageDeclaration, w io.Writer) error {
	if pkg.Name == "" {
		return fmt.Errorf("empty package name")
	}
	file := jen.NewFile(pkg.Name)
	file.PackageComment("file generated. DO NOT EDIT.")

	for _, typ := range pkg.Types {
		err := generateType(typ, file)
		if err != nil {
			return err
		}
	}

	return file.Render(w)
}

func generateInterface(itf *idl.InterfaceType, f *jen.File) error {
	panic("not yet implemented")
}

func generateType(typ signature.Type, f *jen.File) error {

	itf, ok := typ.(*idl.InterfaceType)
	if ok {
		return generateInterface(itf, f)
	}
	typ.TypeDeclaration(f)
	return nil
}
