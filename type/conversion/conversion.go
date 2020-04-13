package conversion

import (
	"reflect"
	"fmt"

	"github.com/lugu/qiloop/type/encoding"
)

// IsConvertibleInto returns true if the `from` type can be
// converted into the `into` type.
func IsConvertibleInto(from, into reflect.Type) bool {
	// TODO: see conversion rules.
	panic("not yet implemented")
}

func AsInt64(you interface{}) (int64, bool) {
	switch y := you.(type) {
	case int:
		return int64(y), true
	case *int:
		return int64(*y), true
	case uint:
		return int64(y), true
	case *uint:
		return int64(*y), true
	case int8:
		return int64(y), true
	case *int8:
		return int64(*y), true
	case uint8:
		return int64(y), true
	case *uint8:
		return int64(*y), true
	case int16:
		return int64(y), true
	case *int16:
		return int64(*y), true
	case uint16:
		return int64(y), true
	case *uint16:
		return int64(*y), true
	case int32:
		return int64(y), true
	case *int32:
		return int64(*y), true
	case uint32:
		return int64(y), true
	case *uint32:
		return int64(*y), true
	case int64:
		return int64(y), true
	case *int64:
		return int64(*y), true
	case uint64:
		return int64(y), true
	case *uint64:
		return int64(*y), true
	default:
		return 0, false
	}
}

// ConvertFrom applies QiMessaging compatibility rules to convert
// `you` into `me`. Me can be populated with default values.
func ConvertFrom(me, you interface{}) error {

	errConversion := fmt.Errorf("Failed to convert %s into %s",
	reflect.TypeOf(you).Name(), reflect.TypeOf(me).Name())

	switch m := me.(type) {
	case *string:
		switch y := you.(type) {
		case string:
			*m = y
			return nil
		case *string:
			*m = *y
			return nil
		default:
			return errConversion
		}
	case *bool:
		switch y := you.(type) {
		case bool:
			*m = y
			return nil
		case *bool:
			*m = *y
			return nil
		default:
			return errConversion
		}
	case *int:
		if in, ok := AsInt64(you); ok {
			*m = int(in)
			return nil
		}
	case *uint:
		if in, ok := AsInt64(you); ok {
			*m = uint(in)
			return nil
		}
	case *int16:
		if in, ok := AsInt64(you); ok {
			*m = int16(in)
			return nil
		}
	case *uint16:
		if in, ok := AsInt64(you); ok {
			*m = uint16(in)
			return nil
		}
	case *int32:
		if in, ok := AsInt64(you); ok {
			*m = int32(in)
			return nil
		}
	case *uint32:
		if in, ok := AsInt64(you); ok {
			*m = uint32(in)
			return nil
		}
	case *int64:
		if in, ok := AsInt64(you); ok {
			*m = int64(in)
			return nil
		}
	case *uint64:
		if in, ok := AsInt64(you); ok {
			*m = uint64(in)
			return nil
		}
	case *float32:
		switch y := you.(type) {
		case float32:
			*m = float32(y)
			return nil
		case *float32:
			*m = float32(*y)
			return nil
		case float64:
			*m = float32(y)
			return nil
		case *float64:
			*m = float32(*y)
			return nil
		default:
			return errConversion
		}
	case *float64:
		switch y := you.(type) {
		case float32:
			*m = float64(y)
			return nil
		case *float32:
			*m = float64(*y)
			return nil
		case float64:
			*m = float64(y)
			return nil
		case *float64:
			*m = float64(*y)
			return nil
		default:
			return errConversion
		}
	}
	return errConversion
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
