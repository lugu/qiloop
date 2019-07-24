package signature

import (
	"bytes"
	"fmt"
	"io"

	"github.com/lugu/qiloop/type/basic"
)

// TypeReader reads from r the bytes representing the type
type TypeReader interface {
	Read(r io.Reader) ([]byte, error)
}

// constReader is a Reader which always read a constant size.
type constReader int

func (c constReader) Read(r io.Reader) ([]byte, error) {
	data := make([]byte, int(c))
	err := basic.ReadN(r, data, int(c))
	if err != nil {
		return nil, err
	}
	return data, nil
}

type stringReader struct{}

func (v stringReader) Read(r io.Reader) ([]byte, error) {
	str, err := basic.ReadString(r)
	var buf bytes.Buffer
	err = basic.WriteString(str, &buf)
	return buf.Bytes(), err
}

// UnknownReader is a TypeReader which returns an error.
type UnknownReader string

func (v UnknownReader) Read(r io.Reader) ([]byte, error) {
	return nil, fmt.Errorf("Unknown type '%v'", v)
}

type objectReader struct{}

func (v objectReader) Read(r io.Reader) ([]byte, error) {
	// TODO: hand make reader, save metaobject reader.
	panic("not yet imple")
}

type valueReader struct{}

func (v valueReader) Read(r io.Reader) ([]byte, error) {
	sig, err := basic.ReadString(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read signature: %s", err)
	}
	reader, err := MakeReader(sig)
	if err != nil {
		return nil, err
	}
	data, err := reader.Read(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read value: %s", err)
	}
	return append([]byte(sig), data...), err
}

type varReader struct {
	reader TypeReader
}

func (v varReader) Read(r io.Reader) ([]byte, error) {
	size, err := basic.ReadUint32(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read size: %s", err)
	}
	if int(size) < 0 {
		return nil, fmt.Errorf("invalid size: %d", size)
	}
	var buf bytes.Buffer
	err = basic.WriteUint32(size, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write size %d: %s",
			size, err)
	}
	for i := 0; i < int(size); i++ {
		data, err := v.reader.Read(r)
		if err != nil {
			return nil, fmt.Errorf("failed to read %d/%d: %s",
				i+1, size, err)
		}
		n, err := buf.Write(data)
		if err != nil {
			return nil, fmt.Errorf("failed to copy %d/%d: %s",
				i, size, err)
		}
		if n != len(data) {
			return nil, fmt.Errorf("failed to copy %d/%d", i, size)
		}
	}
	return buf.Bytes(), nil
}

type tupleReader map[string]TypeReader

func (v tupleReader) Read(r io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	for name, reader := range v {
		data, err := reader.Read(r)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %s",
				name, err)
		}
		n, err := buf.Write(data)
		if err != nil {
			return nil, fmt.Errorf("failed to write %s: %s",
				name, err)
		}
		if n != len(data) {
			return nil, fmt.Errorf("failed to write %d/%d",
				n, len(data))
		}
	}
	return buf.Bytes(), nil
}

// MakeReader parse the signature and returns its associated Reader.
func MakeReader(sig string) (TypeReader, error) {
	t, err := Parse(sig)
	if err != nil {
		return nil, err
	}
	return t.Reader(), nil
}
