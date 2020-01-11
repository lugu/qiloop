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
	A int16
	B int32
	C int64
	D uint16
	E uint32
	F uint64
	G string
	H bool
	I float32
	J float64
}

// readInfoBaseLine unmarshalls InfoBaseLine
func readInfoBaseLine(r io.Reader) (s InfoBaseLine, err error) {
	if s.A, err = basic.ReadInt16(r); err != nil {
		return s, fmt.Errorf("read A field: %s", err)
	}
	if s.B, err = basic.ReadInt32(r); err != nil {
		return s, fmt.Errorf("read B field: %s", err)
	}
	if s.C, err = basic.ReadInt64(r); err != nil {
		return s, fmt.Errorf("read C field: %s", err)
	}
	if s.D, err = basic.ReadUint16(r); err != nil {
		return s, fmt.Errorf("read D field: %s", err)
	}
	if s.E, err = basic.ReadUint32(r); err != nil {
		return s, fmt.Errorf("read E field: %s", err)
	}
	if s.F, err = basic.ReadUint64(r); err != nil {
		return s, fmt.Errorf("read F field: %s", err)
	}
	if s.G, err = basic.ReadString(r); err != nil {
		return s, fmt.Errorf("read G field: %s", err)
	}
	if s.H, err = basic.ReadBool(r); err != nil {
		return s, fmt.Errorf("read H field: %s", err)
	}
	if s.I, err = basic.ReadFloat32(r); err != nil {
		return s, fmt.Errorf("read I field: %s", err)
	}
	if s.J, err = basic.ReadFloat64(r); err != nil {
		return s, fmt.Errorf("read J field: %s", err)
	}
	return s, nil
}

// writeInfoBaseLine marshalls InfoBaseLine
func writeInfoBaseLine(s InfoBaseLine, w io.Writer) (err error) {
	if err := basic.WriteInt16(s.A, w); err != nil {
		return fmt.Errorf("write A field: %s", err)
	}
	if err := basic.WriteInt32(s.B, w); err != nil {
		return fmt.Errorf("write B field: %s", err)
	}
	if err := basic.WriteInt64(s.C, w); err != nil {
		return fmt.Errorf("write C field: %s", err)
	}
	if err := basic.WriteUint16(s.D, w); err != nil {
		return fmt.Errorf("write D field: %s", err)
	}
	if err := basic.WriteUint32(s.E, w); err != nil {
		return fmt.Errorf("write E field: %s", err)
	}
	if err := basic.WriteUint64(s.F, w); err != nil {
		return fmt.Errorf("write F field: %s", err)
	}
	if err := basic.WriteString(s.G, w); err != nil {
		return fmt.Errorf("write G field: %s", err)
	}
	if err := basic.WriteBool(s.H, w); err != nil {
		return fmt.Errorf("write H field: %s", err)
	}
	if err := basic.WriteFloat32(s.I, w); err != nil {
		return fmt.Errorf("write I field: %s", err)
	}
	if err := basic.WriteFloat64(s.J, w); err != nil {
		return fmt.Errorf("write J field: %s", err)
	}
	return nil
}
