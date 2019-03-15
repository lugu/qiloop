//go:generate go get github.com/lugu/qiloop/cmd/qiloop
//go:generate $GOPATH/bin/qiloop stub --idl object.idl --output object_stub_gen.go --path github.com/lugu/qiloop/bus/server/generic

package generic
