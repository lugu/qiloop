package encoding

import (
	"fmt"
	"io"
	"reflect"

	"github.com/lugu/qiloop/type/basic"
)

const (
	listValueMaxSize = 4096
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
	case reflect.Bool:
		return basic.WriteBool(v.Bool(), q.w)
	case reflect.String:
		return basic.WriteString(v.String(), q.w)
	case reflect.Int16:
		return basic.WriteInt16(int16(v.Int()), q.w)
	case reflect.Int32:
		return basic.WriteInt32(int32(v.Int()), q.w)
	case reflect.Int64, reflect.Int:
		return basic.WriteInt64(v.Int(), q.w)
	case reflect.Uint16:
		return basic.WriteUint16(uint16(v.Uint()), q.w)
	case reflect.Uint32:
		return basic.WriteUint32(uint32(v.Uint()), q.w)
	case reflect.Uint64, reflect.Uint:
		return basic.WriteUint64(v.Uint(), q.w)
	case reflect.Float32:
		return basic.WriteFloat32(float32(v.Float()), q.w)
	case reflect.Float64:
		return basic.WriteFloat64(float64(v.Float()), q.w)
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
				if err := q.value(v); err != nil {
					return err
				}
			}
		}
	case reflect.Slice:
		l := v.Len()
		err := basic.WriteInt32(int32(l), q.w)
		if err != nil {
			return fmt.Errorf("slice size write: %w", err)
		}
		for i := 0; i < l; i++ {
			if err = q.value(v.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		keys := v.MapKeys()
		//fmt.Printf("%#v\n", keys)
		err := basic.WriteInt32(int32(len(keys)), q.w)
		if err != nil {
			return fmt.Errorf("map size write: %w", err)
		}
		for _, k := range keys {
			if err = q.value(k); err != nil {
				return err
			}
			if err = q.value(v.MapIndex(k)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (q qiEncoder) Encode(x interface{}) error {
	switch v := x.(type) {
	case string:
		return basic.WriteString(v, q.w)
	case bool:
		return basic.WriteBool(v, q.w)
	case int:
		return basic.WriteInt64(int64(v), q.w)
	case uint:
		return basic.WriteUint64(uint64(v), q.w)
	case uint16:
		return basic.WriteUint16(v, q.w)
	case uint32:
		return basic.WriteUint32(v, q.w)
	case uint64:
		return basic.WriteUint64(v, q.w)
	case int16:
		return basic.WriteInt16(v, q.w)
	case int32:
		return basic.WriteInt32(v, q.w)
	case int64:
		return basic.WriteInt64(v, q.w)
	case float32:
		return basic.WriteFloat32(v, q.w)
	case float64:
		return basic.WriteFloat64(v, q.w)
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
		if l > listValueMaxSize {
			return fmt.Errorf("list too long: %d", l)
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
		if l > listValueMaxSize {
			return fmt.Errorf("list too long: %d", l)
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
		if l > listValueMaxSize {
			return fmt.Errorf("map too long: %d", l)
		}
		v.Set(reflect.MakeMapWithSize(v.Elem().Type(), l))
		v = v.Elem()
	}
	if v.Kind() != reflect.Map {
		return fmt.Errorf("not a slice: %v", v)
	}
	if v.IsNil() {
		if l > listValueMaxSize {
			return fmt.Errorf("map too long: %d", l)
		}
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
	case reflect.Bool:
		b, err := basic.ReadBool(q.r)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.String:
		s, err := basic.ReadString(q.r)
		if err != nil {
			return err
		}
		v.SetString(s)
	case reflect.Int16:
		i, err := basic.ReadInt16(q.r)
		if err != nil {
			return err
		}
		v.SetInt(int64(i))
	case reflect.Int32:
		i, err := basic.ReadInt32(q.r)
		if err != nil {
			return err
		}
		v.SetInt(int64(i))
	case reflect.Int64, reflect.Int:
		i, err := basic.ReadInt64(q.r)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint16:
		i, err := basic.ReadUint16(q.r)
		if err != nil {
			return err
		}
		v.SetUint(uint64(i))
	case reflect.Uint32:
		i, err := basic.ReadUint32(q.r)
		if err != nil {
			return err
		}
		v.SetUint(uint64(i))
	case reflect.Uint64, reflect.Uint:
		i, err := basic.ReadUint64(q.r)
		if err != nil {
			return err
		}
		v.SetUint(uint64(i))
	case reflect.Float32:
		f, err := basic.ReadFloat32(q.r)
		if err != nil {
			return err
		}
		v.SetFloat(float64(f))
	case reflect.Float64:
		f, err := basic.ReadFloat64(q.r)
		if err != nil {
			return err
		}
		v.SetFloat(f)
	}
	return nil
}

func (q qiDecoder) Decode(x interface{}) (err error) {
	switch v := x.(type) {
	case *string:
		*v, err = basic.ReadString(q.r)
		return err
	case *bool:
		*v, err = basic.ReadBool(q.r)
		return err
	case *int:
		var tmp int64
		tmp, err = basic.ReadInt64(q.r)
		*v = int(tmp)
		return err
	case *uint:
		var tmp uint64
		tmp, err = basic.ReadUint64(q.r)
		*v = uint(tmp)
		return err
	case *uint16:
		*v, err = basic.ReadUint16(q.r)
		return err
	case *uint32:
		*v, err = basic.ReadUint32(q.r)
		return err
	case *uint64:
		*v, err = basic.ReadUint64(q.r)
		return err
	case *int16:
		*v, err = basic.ReadInt16(q.r)
		return err
	case *int32:
		*v, err = basic.ReadInt32(q.r)
		return err
	case *int64:
		*v, err = basic.ReadInt64(q.r)
		return err
	case *float32:
		*v, err = basic.ReadFloat32(q.r)
		return err
	case *float64:
		*v, err = basic.ReadFloat64(q.r)
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
