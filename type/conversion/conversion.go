package conversion

import (
	"reflect"
	"strings"
	"fmt"

	"github.com/lugu/qiloop/type/encoding"
)

// IsConvertibleInto returns true if the `from` type can be
// converted into the `into` type.
func IsConvertibleInto(from, into reflect.Type) bool {
	// TODO: see conversion rules.
	panic("not yet implemented")
}

func AsInt64(w reflect.Value) (int64, bool) {
	switch w.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		return w.Int(), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		return int64(w.Uint()), true
	default:
		return 0, false
	}
}

func convertSlice(v, w reflect.Value) error {
	if w.Kind() != reflect.Slice {
		return fmt.Errorf("Failed to convert slice %v into %v",
			v.Type(), w.Type())
	}

	l := w.Len()
	if v.Kind() == reflect.Ptr && v.IsNil() {
		if !v.CanSet() {
			return fmt.Errorf("cannot set slice: %v", v)
		}
		v.Set(reflect.MakeSlice(v.Elem().Type(), l, l))
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("not a slice: %v", v)
	}
	if v.Cap() < l {
		if v.CanSet() == false {
			return fmt.Errorf("slice capacity too short, cannot set : %d", v.Cap())
		}
		v.Set(reflect.MakeSlice(v.Type(), l, l))
	}
	v.SetLen(l)
	for i := 0; i < l; i++ {
		err := convertFrom(v.Index(i), w.Index(i))
		if err != nil {
			return fmt.Errorf("cannot convert slice index %d, %w", i, err)
		}
	}
	return nil
}

func convertMap(v, w reflect.Value) error {
	if w.Kind() != reflect.Map {
		return fmt.Errorf("Failed to convert map %v into %v",
			v.Type(), w.Type())
	}
	l := w.Len()
	if v.Kind() == reflect.Ptr && v.IsNil() {
		if !v.CanSet() {
			return fmt.Errorf("cannot set slice: %v", v)
		}
		v.Set(reflect.MakeMapWithSize(v.Elem().Type(), l))
		v = v.Elem()
	}
	if v.Kind() != reflect.Map {
		return fmt.Errorf("not a map: %v", v)
	}
	if v.IsNil() {
		v.Set(reflect.MakeMapWithSize(v.Type(), l))
	}
	for _, k := range w.MapKeys() {
		key := reflect.New(v.Type().Key())
		err := convertFrom(key, k)
		if err != nil {
			return fmt.Errorf("cannot convert map key %w", err)
		}

		el := reflect.New(v.Type().Elem())
		err = convertFrom(key, w.MapIndex(k))
		if err != nil {
			return fmt.Errorf("cannot convert map value %w", err)
		}
		v.SetMapIndex(key.Elem(), el.Elem())
	}
	return nil
}

func convertStruct(v, w reflect.Value) error {
	if w.Kind() != reflect.Struct {
		return fmt.Errorf("Failed to convert struct %v into %v",
			v.Type(), w.Type())
	}
	if v.Kind() == reflect.Ptr && v.IsNil() {
		if !v.CanSet() {
			return fmt.Errorf("cannot set struct ptr: %v", v)
		}
		v.Set(reflect.New(v.Elem().Type()))
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		name := strings.ToLower(v.Type().Field(i).Name)
		for j := 0; j < w.NumField(); j++ {
			if name == strings.ToLower(w.Type().Field(j).Name) {
				err := convertFrom( v.Field(i), w.Field(j))
				if err != nil {
					return fmt.Errorf("failed to convert filed %s of %v: %w",
						name, v.Type(), err)
				}
				break
			}
		}

	}
	return nil
}

func convertFrom(v, w reflect.Value) error {

	errConversion := fmt.Errorf("Failed to convert %v into %v",
		v.Type(), w.Type())

	if w.Kind() == reflect.Ptr {
		return convertFrom(v, w.Elem())
	}

	switch v.Kind() {
	case reflect.Bool:
		if w.Kind() == reflect.Bool {
			v.SetBool(w.Bool())
			return nil
		}
	case reflect.String:
		if w.Kind() == reflect.String {
			v.SetString(w.String())
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64:
		if i, ok := AsInt64(w); ok {
			v.SetInt(i)
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		if i, ok := AsInt64(w); ok {
			v.SetUint(uint64(i))
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if w.Kind() == reflect.Float32 ||
			w.Kind() == reflect.Float64 {
			v.SetFloat(w.Float())
			return nil
		}
	case reflect.Ptr:
		v = v.Elem()
		switch v.Kind() {
		case reflect.Slice:
			return convertSlice(v, w)
		case reflect.Map:
			return convertMap(v, w)
		case reflect.Struct:
			return convertStruct(v, w)
		default:
			return convertFrom(v, w)
		}
	case reflect.Slice:
		return convertSlice(v, w)
	case reflect.Map:
		return convertMap(v, w)
	case reflect.Struct:
		return convertStruct(v, w)
	}

	return errConversion
}

// ConvertFrom applies QiMessaging compatibility rules to convert
// `you` into `me`. Me can be populated with default values.
func ConvertFrom(me, you interface{}) error {
	return convertFrom(reflect.ValueOf(me),reflect.ValueOf(you))
}

// DecodeFrom decodes x from an encoded value of tye `typ`.
func DecodeFrom(d encoding.Decoder, x interface{}, typ reflect.Type) error {
	from := reflect.New(typ)
	if err := d.Decode(from.Interface()); err != nil {
		return err
	}
	return convertFrom(reflect.ValueOf(x), from)
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
