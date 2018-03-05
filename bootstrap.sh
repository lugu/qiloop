#!/bin/sh

go generate github.com/lugu/qiloop/meta/stage1
go generate github.com/lugu/qiloop/meta/stage2
go generate github.com/lugu/qiloop/meta/stage3
go generate github.com/lugu/qiloop/services
