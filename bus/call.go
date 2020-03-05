package bus

import (
	"github.com/lugu/qiloop/type/encoding"
)


// Params represents a tuple of values to be sent
type Params struct {
	sig  string
	args []interface{}
}

// NewParams returns a
func NewParams(sig string, args... interface{}) Params {
	return Params{
		sig: sig,
		args: args,
	}
}

func (p *Params) Signature() string {
	return p.sig
}

func (p *Params) Write(e encoding.Encoder) error {
	for _, arg := range p.args {
		if err := e.Encode(arg); err != nil {
			return err
		}
	}
	return nil
}

type Response struct {
	sig  string
	resp interface{}
}

func NewResponse(signature string, instance interface{}) Response {
	return Response{
		sig: signature,
		resp: instance,
	}
}

// Signature returns the signature of the value.
func (o *Response) Signature() string {
	return o.sig
}

func (o *Response) Read(d encoding.Decoder) error {
	return d.Decode(o.resp)
}
