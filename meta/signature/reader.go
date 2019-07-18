package signature

import (
	"fmt"
	"io"
)

// TypeReader reads from r the bytes representing the type
type TypeReader interface {
	Read(r io.Reader) ([]byte, error)
}

// ConstReader is a Reader which always read a constant size.
type ConstReader int

func (c ConstReader) Read(r io.Reader) ([]byte, error) {
	data := make([]byte, int(c))
	n, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	if n != int(c) {
		return nil, fmt.Errorf("type requires %d bytes, read %d byte",
			int(c), n)
	}
	return data, nil
}

// MakeReader parse the signature and returns its associated Reader.
func MakeReader(sig string) (TypeReader, error) {
	panic("not yet implemented")
}
