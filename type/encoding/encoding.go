package encoding

import (
	"fmt"
	"io"
	"reflect"

	"github.com/lugu/qiloop/type/basic"
)

type BinaryEncoder interface {
	Encode(Encoder) error
}
type Encoder interface {
	Encode(interface{}) error
}

type qiEncoder struct {
	w io.Writer
}

func (q qiEncoder) value(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr:
		return q.value(v.Elem())
	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			// Note: Calling v.CanSet() below is an optimization.
			// It would be sufficient to check the field name,
			// but creating the StructField info for each field is
			// costly (run "go test -bench=ReadStruct" and compare
			// results when making changes to this code).
			if v := v.Field(i); v.CanSet() || t.Field(i).Name != "_" {
				q.value(v)
			}
		}
	case reflect.Slice:
		l := v.Len()
		err := basic.WriteInt32(int32(l), q.w)
		if err != nil {
			return fmt.Errorf("slice size write: %w", err)
		}
		for i := 0; i < l; i++ {
			q.value(v.Index(i))
		}
	case reflect.Map:
		keys := v.MapKeys()
		//fmt.Printf("%#v\n", keys)
		err := basic.WriteInt32(int32(len(keys)), q.w)
		if err != nil {
			return fmt.Errorf("map size write: %w", err)
		}
		for _, k := range keys {
			q.value(k)
			q.value(v.MapIndex(k))
		}
	case reflect.Int16:
		err := basic.WriteInt16(int16(v.Int()), q.w)
		if err != nil {
			return fmt.Errorf("failed to write int16: %w", err)
		}
	case reflect.Int32:
		err := basic.WriteInt32(int32(v.Int()), q.w)
		if err != nil {
			return fmt.Errorf("failed to write int32: %w", err)
		}
	case reflect.Int64:
		err := basic.WriteInt64(int64(v.Int()), q.w)
		if err != nil {
			return fmt.Errorf("failed to write int64: %w", err)
		}
	}
	return nil
}

func (q qiEncoder) Encode(x interface{}) error {
	switch v := x.(type) {
	case int16:
		return basic.WriteInt16(v, q.w)
	case int32:
		return basic.WriteInt32(v, q.w)
	case int64:
		return basic.WriteInt64(v, q.w)
	case BinaryEncoder:
		return v.Encode(q)
	}

	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Ptr: // ok
	case reflect.Slice: // ok
	case reflect.Map: // ok
	case reflect.Struct: // ok
	default:
		return fmt.Errorf("can only read from pointer, map or slice, not kind: %d", v.Kind())
	}

	// Fallback to reflect-based encoding.
	return q.value(v)
}

func NewEncoder(permission map[string]string, w io.Writer) Encoder {
	return qiEncoder{w}
}

type BinaryDecoder interface {
	Decode(Decoder) error
}
type Decoder interface {
	Decode(interface{}) error
}

type qiDecoder struct {
	r io.Reader
}

func (q qiDecoder) sliceValue(v reflect.Value) error {
	length, err := basic.ReadInt32(q.r)
	if err != nil {
		return fmt.Errorf("failed to read vector size: %w", err)
	}
	l := int(length)
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
		q.value(v.Index(i))
	}
	return nil
}

func (q qiDecoder) mapValue(v reflect.Value) error {
	length, err := basic.ReadInt32(q.r)
	if err != nil {
		return fmt.Errorf("failed to read vector size: %w", err)
	}
	l := int(length)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		if !v.CanSet() {
			return fmt.Errorf("cannot set slice: %v", v)
		}
		v.Set(reflect.MakeMapWithSize(v.Elem().Type(), l))
		v = v.Elem()
	}
	if v.Kind() != reflect.Map {
		return fmt.Errorf("not a slice: %v", v)
	}
	if v.IsNil() {
		v.Set(reflect.MakeMapWithSize(v.Type(), l))
	}
	for i := 0; i < l; i++ {
		key := reflect.New(v.Type().Key())
		err := q.value(key)
		if err != nil {
			return fmt.Errorf("read map key failed: %w", err)
		}
		el := reflect.New(v.Type().Elem())
		err = q.value(el)
		if err != nil {
			return fmt.Errorf("read map element failed: %w", err)
		}
		v.SetMapIndex(key.Elem(), el.Elem())
	}
	return nil
}

func (q qiDecoder) value(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
		if v.Kind() == reflect.Slice {
			return q.sliceValue(v)
		} else if v.Kind() == reflect.Map {
			return q.mapValue(v)
		}
		return q.value(v)
	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			// Note: Calling v.CanSet() below is an optimization.
			// It would be sufficient to check the field name,
			// but creating the StructField info for each field is
			// costly (run "go test -bench=ReadStruct" and compare
			// results when making changes to this code).
			if v := v.Field(i); v.CanSet() || t.Field(i).Name != "_" {
				q.value(v)
			}
		}
	case reflect.Slice:
		return q.sliceValue(v)
	case reflect.Map:
		return q.mapValue(v)
	case reflect.Int16:
		i, err := basic.ReadInt16(q.r)
		if err != nil {
			return fmt.Errorf("failed to read int16: %w", err)
		}
		v.SetInt(int64(i))
	case reflect.Int32:
		i, err := basic.ReadInt32(q.r)
		if err != nil {
			return fmt.Errorf("failed to read int32: %w", err)
		}
		v.SetInt(int64(i))
	case reflect.Int64:
		i, err := basic.ReadInt64(q.r)
		if err != nil {
			return fmt.Errorf("failed to read int64: %w", err)
		}
		v.SetInt(i)
	}
	return nil
}

func (q qiDecoder) Decode(x interface{}) (err error) {
	switch v := x.(type) {
	case *int16:
		*v, err = basic.ReadInt16(q.r)
		return err
	case *int32:
		*v, err = basic.ReadInt32(q.r)
		return err
	case *int64:
		*v, err = basic.ReadInt64(q.r)
		return err
	case BinaryDecoder:
		return v.Decode(q)
	}

	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Ptr: // ok
	case reflect.Slice: // ok
	case reflect.Map: // ok
	default:
		return fmt.Errorf("can only read from pointer, map or slice, not kind: %d", v.Kind())
	}

	// Fallback to reflect-based decoding.
	return q.value(v)
}

func NewDecoder(permissions map[string]string, r io.Reader) Decoder {
	return qiDecoder{r}
}
