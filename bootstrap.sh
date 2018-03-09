#!/bin/sh

set -e

rm -f \
    meta/stage1/metaobject.go \
    meta/stage2/services.go   \
    meta/stage3/services.go   \
    services/services.go      \
    object/metaobject.go      \

go get github.com/lugu/qiloop/session
go get github.com/lugu/qiloop/meta/signature
go get github.com/lugu/qiloop/meta/stage1/...
go generate github.com/lugu/qiloop/meta/stage1
go get github.com/lugu/qiloop/object/...
go generate github.com/lugu/qiloop/object
go get github.com/lugu/qiloop/meta/proxy
go get github.com/lugu/qiloop/meta/stage2/...
go generate github.com/lugu/qiloop/meta/stage2
go get github.com/lugu/qiloop/meta/stage3/...
go generate github.com/lugu/qiloop/meta/stage3
go generate github.com/lugu/qiloop/services
