package conversion

import (
	"reflect"

	"github.com/lugu/qiloop/type/encoding"
)

// Parse returns a Go type which represents the signature.
func Parse(signature string) (reflect.Type, error) {
	// TODO: update signature package to provide a reflect.Type.
	panic("not yet implemented")
}

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

// IsSignatureCompatible returns true if a type described by the
// signature `from` can be encoded into something compatible a
// type described by the signature `into`.
func IsSignatureCompatible(from, into string) (bool, error) {
	fromT, err := Parse(from)
	if err != nil {
		return false, err
	}
	intoT, err := Parse(into)
	if err != nil {
		return false, err
	}
	return IsConvertibleInto(fromT, intoT), nil
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
