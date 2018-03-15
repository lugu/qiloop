#!/usr/bin/env bash

set -e

rm -f \
    meta/stage2/interfaces.go \
    meta/stage2/services.go   \
    meta/stage3/interfaces.go \
    meta/stage3/services.go   \
    bus/services/interfaces.go    \
    bus/services/services.go      \
    type/object/metaobject.go      \

go get -d github.com/lugu/qiloop/...
go generate github.com/lugu/qiloop/type/object
go generate github.com/lugu/qiloop/meta/stage2
go generate github.com/lugu/qiloop/meta/stage3
go generate github.com/lugu/qiloop/bus/services
go get github.com/lugu/qiloop/cmd/info

go test github.com/lugu/qiloop/... -race -coverprofile=coverage.txt -covermode=atomic
$HOME/go/bin/info | head
