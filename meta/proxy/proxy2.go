package proxy

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/signature"
	"io"
)

// GeneratePackage generate the proxy for the package declaration.
func GeneratePackage(w io.Writer, pkg *idl.PackageDeclaration) error {

	if pkg.Name == "" {
		return fmt.Errorf("empty package name")
	}
	file := jen.NewFile(pkg.Name)
	msg := "Package " + pkg.Name + " contains a generated proxy"
	file.PackageComment(msg)
	file.PackageComment("File generated. DO NOT EDIT.")

	set := signature.NewTypeSet()
	for _, typ := range pkg.Types {
		err := generateType(file, set, typ)
		if err != nil {
			return err
		}
	}

	return file.Render(w)
}

func generateType(file *jen.File, set *signature.TypeSet, typ signature.Type) error {

	itf, ok := typ.(*idl.InterfaceType)
	if ok {
		metaObj := itf.MetaObject()
		serviceName := util.CleanName(metaObj.Description)
		generateObjectInterface(metaObj, serviceName, set, file)

		err := generateProxyObject(metaObj, serviceName, set, file)
		if err != nil {
			return fmt.Errorf("failed to render %s: %s",
				metaObj.Description, err)
		}
	}
	set.Declare(file)
	return nil
}
