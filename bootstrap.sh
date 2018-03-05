#!/bin/sh

set -e

go get github.com/lugu/qiloop/net
go get github.com/lugu/qiloop/meta/signature
go get github.com/lugu/qiloop/meta/stage1/...
go generate github.com/lugu/qiloop/meta/stage1
go get github.com/lugu/qiloop/meta/proxy
go get github.com/lugu/qiloop/meta/stage2/...
go generate github.com/lugu/qiloop/meta/stage2
go get github.com/lugu/qiloop/meta/stage3/...
go generate github.com/lugu/qiloop/meta/stage3
go get github.com/lugu/qiloop/meta/stage4/...
go generate github.com/lugu/qiloop/meta/stage4
