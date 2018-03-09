#!/usr/bin/env bash

set -e

rm -f \
    meta/stage1/metaobject.go \
    meta/stage2/services.go   \
    meta/stage3/services.go   \
    services/services.go      \
    object/metaobject.go      \

go get -d github.com/lugu/qiloop/...
go generate github.com/lugu/qiloop/meta/stage1
go generate github.com/lugu/qiloop/object
go generate github.com/lugu/qiloop/meta/stage2
go generate github.com/lugu/qiloop/meta/stage3
go generate github.com/lugu/qiloop/services
