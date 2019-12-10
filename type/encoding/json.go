package encoding

import (
	"encoding/json"
	"io"

	"github.com/lugu/qiloop/type/basic"
)

type jsonEncoder struct {
	w io.Writer
}

func (j jsonEncoder) Encode(x interface{}) error {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return err
	}
	return basic.WriteN(j.w, b, len(b))
}

func NewJSONEncoder(w io.Writer) Encoder {
	return jsonEncoder{w}
}

type jsonDecoder struct {
	b []byte
}

func (j jsonDecoder) Decode(x interface{}) error {
	return json.Unmarshal(j.b, x)
}

func NewJSONDecoder(d []byte) Decoder {
	return jsonDecoder{d}
}
