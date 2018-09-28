package fuzz_test

import (
	"github.com/lugu/qiloop/bus/cmd/fuzz"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestFuzz(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic with %s", r)
		}
	}()
	data, err := ioutil.ReadFile(filepath.Join("testdata", "cap-auth-failure.bin"))
	if err != nil {
		t.Errorf("cannot open test data %s", err)
	}
	fuzz.Fuzz(data)
}
