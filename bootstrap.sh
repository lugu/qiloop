#!/usr/bin/env bash

set -e

rm -f \
    meta/stage2/interfaces.go \
    meta/stage2/services.go   \
    meta/stage3/interfaces.go \
    meta/stage3/services.go   \
    services/interfaces.go    \
    services/services.go      \
    object/metaobject.go      \

go get -d github.com/lugu/qiloop/...
go generate github.com/lugu/qiloop/object
go generate github.com/lugu/qiloop/meta/stage2
go generate github.com/lugu/qiloop/meta/stage3
go generate github.com/lugu/qiloop/services
go get github.com/lugu/qiloop/cmd/info

go test github.com/lugu/qiloop/... -race -coverprofile=coverage.txt -covermode=atomic
$HOME/go/bin/info | head
