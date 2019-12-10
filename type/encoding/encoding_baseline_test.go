package encoding_test

// Baseline encoder, made with:
// go run ../../meta/cmd/stub --idl test.idl --output encoding_baseline_test.go

import (
	"fmt"
	"io"

	basic "github.com/lugu/qiloop/type/basic"
)

// InfoBaseLine is serializable
type InfoBaseLine struct {
	A int32
	B int16
	C int64
}

// readInfoBaseLine unmarshalls InfoBaseLine
func readInfoBaseLine(r io.Reader) (s InfoBaseLine, err error) {
	if s.A, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("read A field: %s", err)
	}
	if s.B, err = basic.ReadInt16(r); err != nil {
		return s, fmt.Errorf("read B field: %s", err)
	}
	if s.C, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("read C field: %s", err)
	}
	return s, nil
}

// writeInfoBaseLine marshalls InfoBaseLine
func writeInfoBaseLine(s InfoBaseLine, w io.Writer) (err error) {
	if err := basic.WriteInt32(s.A, w); err != nil {
		return fmt.Errorf("write A field: %s", err)
	}
	if err := basic.WriteInt16(s.B, w); err != nil {
		return fmt.Errorf("write B field: %s", err)
	}
	if err := basic.WriteInt64(s.C, w); err != nil {
		return fmt.Errorf("write C field: %s", err)
	}
	return nil
}
