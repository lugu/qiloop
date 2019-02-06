//go:generate go get github.com/lugu/qiloop/cmd/proxygen
//go:generate $GOPATH/bin/proxygen -idl server.idl -output proxy_gen.go

package client
