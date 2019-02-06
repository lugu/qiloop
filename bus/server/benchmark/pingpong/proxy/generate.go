//go:generate go get github.com/lugu/qiloop/cmd/proxygen
//go:generate $GOPATH/bin/proxygen -idl ping.idl -output proxy_gen.go

package proxy
