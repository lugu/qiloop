package main

import (
	s "github.com/lugu/qiloop/meta/stub"
)

func stub(idlFileName, stubFileName, packageName string) {
	s.GenerateStub(idlFileName, stubFileName, packageName)
}
