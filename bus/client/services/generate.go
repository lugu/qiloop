//go:generate go get github.com/lugu/qiloop/cmd/proxygen
//go:generate $GOPATH/bin/proxygen -idl services.idl -output proxy_gen.go

package services
