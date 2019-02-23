//go:generate go get github.com/lugu/qiloop/cmd/stub
//go:generate $GOPATH/bin/stub -idl object.idl -output object_stub_gen.go -path github.com/lugu/qiloop/bus/server/generic
//go:generate $GOPATH/bin/stub -idl testdata.idl -output data_test.go -path github.com/lugu/qiloop/bus/server/generic

package generic
