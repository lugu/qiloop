package proxy

import (
	"fmt"
	"github.com/dave/jennifer/jen"
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
	file.HeaderComment(msg)
	file.HeaderComment("File generated. DO NOT EDIT.")

	GenerateNewServices(file)

	set := signature.NewTypeSet()
	for _, typ := range pkg.Types {
		typ.RegisterTo(set)
	}

	set.Declare(file)
	return file.Render(w)
}

func GenerateNewServices(file *jen.File) {
	file.Comment("Constructor gives access to remote services")
	file.Type().Id(
		"Constructor",
	).Struct(
		jen.Id("session").Qual("github.com/lugu/qiloop/bus", "Session"),
	)
	file.Comment("Services gives access to the services constructor")
	file.Func().Id(
		"Services",
	).Params(
		jen.Id("s").Qual("github.com/lugu/qiloop/bus", "Session"),
	).Id(
		"Constructor",
	).Block(
		jen.Id(`return Constructor{ session: s, }`),
	)
}
