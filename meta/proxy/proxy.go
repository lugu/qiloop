package proxy

import (
	"fmt"
	"io"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/lugu/qiloop/meta/idl"
	"github.com/lugu/qiloop/meta/signature"
)

// GeneratePackage generate the proxy for the package declaration.
func GeneratePackage(w io.Writer, packagePath string,
	pkg *idl.PackageDeclaration) error {

	if pkg.Name == "" {
		if strings.Contains(packagePath, "/") {
			pkg.Name = packagePath[strings.LastIndex(packagePath, "/")+1:]
		} else {
			return fmt.Errorf("empty package name")
		}
	}
	file := jen.NewFilePathName(packagePath, pkg.Name)
	msg := "Package " + pkg.Name + " contains a generated proxy"
	file.HeaderComment(msg)
	file.HeaderComment(".")

	GenerateNewServices(file)

	set := signature.NewTypeSet()
	for _, typ := range pkg.Types {
		typ.RegisterTo(set)
	}

	set.Declare(file)
	return file.Render(w)
}

// GenerateNewServices declares the constructor object.
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
