package encoding

import (
	"encoding/gob"
	"io"
)

func NewGobEncoder(w io.Writer) Encoder {
	return gob.NewEncoder(w)
}

func NewGobDecoder(r io.Reader) Decoder {
	return gob.NewDecoder(r)
}
