package conversion

import (
	"reflect"

	"github.com/lugu/qiloop/type/encoding"
)

// IsConvertibleInto returns true if the `from` type can be
// converted into the `into` type.
func IsConvertibleInto(from, into reflect.Type) bool {
	// TODO: see conversion rules.
	panic("not yet implemented")
}

// ConvertFrom applies QiMessaging compatibility rules to convert
// `you` into `me`. Me can be populated with default values.
func ConvertFrom(me, you interface{}) error {
	// TODO: see conversion rules.
	panic("not yet implemented")
}

// DecodeFrom decodes x from an encoded value of tye `typ`.
func DecodeFrom(d encoding.Decoder, x interface{}, typ reflect.Type) error {
	from := reflect.New(typ)
	if err := d.Decode(from.Interface()); err != nil {
		return err
	}
	return ConvertFrom(x, from)
}

// EncodeInto encodes x as if it was of type `typ`.
func EncodeInto(e encoding.Encoder, x interface{}, typ reflect.Type) error {
	// step 1: create an instance of typ called out
	out := reflect.New(typ)
	// step 2: decode this instance with: Decoder.Decode()
	err := ConvertFrom(out, x)
	if err != nil {
		return err
	}
	// step 3: encode out
	if err := e.Encode(out); err != nil {
		return err
	}
	return nil
}
