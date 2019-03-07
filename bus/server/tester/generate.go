//go:generate go get github.com/lugu/qiloop/cmd/proxygen
//go:generate $GOPATH/bin/proxygen -idl tester.idl -output proxy_gen.go
//go:generate go get github.com/lugu/qiloop/cmd/stub
//go:generate $GOPATH/bin/stub -idl tester.idl -output stub_gen.go

package tester
