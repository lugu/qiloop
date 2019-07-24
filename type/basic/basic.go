package basic

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// MaxStringSize the longest string allowed.
const MaxStringSize = uint32(10 * 1024 * 1024)

// ReadN tries and retries to read length bytes from r. Reading length
// with io.EOF is not considered an error. Forwards io.EOF if nothing
// was read.
func ReadN(r io.Reader, buf []byte, length int) error {
	size := 0
	for size < length {
		read, err := r.Read(buf[size:])
		size += read
		if err == nil && read != 0 {
			continue
		} else if err == io.EOF && size == length {
			break
		} else if err == io.EOF && size == 0 {
			return io.EOF
		} else {
			if err == nil {
				err = fmt.Errorf("no progress")
			}
			return fmt.Errorf("read %d instead of %d: %s",
				size, length, err)
		}
	}
	return nil
}

// WriteN tries and retries to write length bytes into w. Writing
// length with io.EOF is not considered an error. Forwards io.EOF if
// nothing was written.
func WriteN(w io.Writer, buf []byte, length int) error {
	size := 0
	for size < length {
		write, err := w.Write(buf[size:])
		size += write
		if err == nil && write != 0 {
			continue
		} else if err == io.EOF && size == length {
			break
		} else if err == io.EOF && size == 0 {
			return io.EOF
		} else {
			if err == nil {
				err = fmt.Errorf("no progress")
			}
			return fmt.Errorf("write %d instead of %d: %s",
				size, length, err)
		}
	}
	return nil
}

// ReadUint8 read an uint8
func ReadUint8(r io.Reader) (uint8, error) {
	buf := []byte{0}
	err := ReadN(r, buf, 1)
	if err != nil {
		return 0, err
	}
	return uint8(buf[0]), nil
}

// WriteUint8 an uint8
func WriteUint8(i uint8, w io.Writer) error {
	buf := []byte{i}
	err := WriteN(w, buf, 1)
	if err != nil {
		return err
	}
	return nil
}

// ReadInt8 reads a little endian int8
func ReadInt8(r io.Reader) (int8, error) {
	i, err := ReadUint8(r)
	return int8(i), err
}

// WriteInt8 writes a little endian int8
func WriteInt8(i int8, w io.Writer) error {
	return WriteUint8(uint8(i), w)
}

// ReadUint16 reads a little endian uint16
func ReadUint16(r io.Reader) (uint16, error) {
	buf := []byte{0, 0}
	err := ReadN(r, buf, 2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf), nil
}

// WriteUint16 writes a little endian uint16
func WriteUint16(i uint16, w io.Writer) error {
	buf := []byte{0, 0}
	binary.LittleEndian.PutUint16(buf, i)
	err := WriteN(w, buf, 2)
	if err != nil {
		return err
	}
	return nil
}

// ReadInt16 reads a little endian int16
func ReadInt16(r io.Reader) (int16, error) {
	i, err := ReadUint16(r)
	return int16(i), err
}

// WriteInt16 writes a little endian int16
func WriteInt16(i int16, w io.Writer) error {
	return WriteUint16(uint16(i), w)
}

// ReadUint32 reads a little endian uint32
func ReadUint32(r io.Reader) (uint32, error) {
	buf := []byte{0, 0, 0, 0}
	err := ReadN(r, buf, 4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf), nil
}

// WriteUint32 writes a little endian uint32
func WriteUint32(i uint32, w io.Writer) error {
	buf := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(buf, i)
	err := WriteN(w, buf, 4)
	if err != nil {
		return err
	}
	return nil
}

// ReadInt32 reads a little endian int32
func ReadInt32(r io.Reader) (int32, error) {
	i, err := ReadUint32(r)
	return int32(i), err
}

// WriteInt32 writes a little endian int32
func WriteInt32(i int32, w io.Writer) error {
	return WriteUint32(uint32(i), w)
}

// ReadUint64 read a little endian uint64
func ReadUint64(r io.Reader) (uint64, error) {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	err := ReadN(r, buf, 8)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf), nil
}

// WriteUint64 writes a little endian uint64
func WriteUint64(i uint64, w io.Writer) error {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(buf, i)
	err := WriteN(w, buf, 8)
	if err != nil {
		return err
	}
	return nil
}

// ReadInt64 reads a little endian int64
func ReadInt64(r io.Reader) (int64, error) {
	i, err := ReadUint64(r)
	return int64(i), err
}

// WriteInt64 writes a little endian int64
func WriteInt64(i int64, w io.Writer) error {
	return WriteUint64(uint64(i), w)
}

// ReadFloat32 read a little endian float32
func ReadFloat32(r io.Reader) (float32, error) {
	buf := []byte{0, 0, 0, 0}
	err := ReadN(r, buf, 4)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint32(buf)
	return math.Float32frombits(bits), nil
}

// WriteFloat32 writes a little endian float32
func WriteFloat32(f float32, w io.Writer) error {
	buf := []byte{0, 0, 0, 0}
	bits := math.Float32bits(f)
	binary.LittleEndian.PutUint32(buf, bits)
	err := WriteN(w, buf, 4)
	if err != nil {
		return err
	}
	return nil
}

// ReadFloat64 read a little endian float64
func ReadFloat64(r io.Reader) (float64, error) {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	err := ReadN(r, buf, 8)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(bits), nil
}

// WriteFloat64 writes a little endian float64
func WriteFloat64(f float64, w io.Writer) error {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	bits := math.Float64bits(f)
	binary.LittleEndian.PutUint64(buf, bits)
	err := WriteN(w, buf, 8)
	if err != nil {
		return err
	}
	return nil
}

// ReadBool read a one byte size binary value
func ReadBool(r io.Reader) (bool, error) {
	u, err := ReadUint8(r)
	if u == 0 {
		return false, err
	}
	return true, err
}

// WriteBool writes a one byte size binary value
func WriteBool(b bool, w io.Writer) error {
	if b {
		return WriteUint8(1, w)
	}
	return WriteUint8(0, w)
}

// ReadString reads a string: first the size of the string is read
// using ReadUint32, then the bytes of the string.
func ReadString(r io.Reader) (string, error) {
	size, err := ReadUint32(r)
	if err != nil {
		return "", fmt.Errorf("read string size: %s", err)
	}
	if size == 0 {
		return "", nil
	}
	if size > MaxStringSize {
		return "", fmt.Errorf("invalid string size: %d", size)
	}
	buf := make([]byte, size)
	err = ReadN(r, buf, int(size))
	if err != nil {
		return "", fmt.Errorf("read string: %s", err)
	}
	return string(buf), nil
}

// WriteString writes a string: first the size of the string is
// written using WriteUint32, then the bytes of the string.
func WriteString(s string, w io.Writer) error {
	if err := WriteUint32(uint32(len(s)), w); err != nil {
		return fmt.Errorf("write string size: %s", err)
	}
	err := WriteN(w, []byte(s), len(s))
	if err != nil {
		return err
	}
	return nil
}
