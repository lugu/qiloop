package basic

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func ReadUint8(r io.Reader) (uint8, error) {
	buf := []byte{0}
	bytes, err := r.Read(buf)
	if err != nil {
		return 0, err
	} else if bytes != 1 {
		return 0, fmt.Errorf("failed to read uint8 (%d instead of 1)", bytes)
	}
	return uint8(buf[0]), nil
}

func WriteUint8(i uint8, w io.Writer) error {
	buf := []byte{i}
	bytes, err := w.Write(buf)
	if err != nil {
		return err
	} else if bytes != 1 {
		return fmt.Errorf("failed to write uint16 (%d instead of 1)", bytes)
	}
	return nil
}

func ReadUint16(r io.Reader) (uint16, error) {
	buf := []byte{0, 0}
	bytes, err := r.Read(buf)
	if err != nil {
		return 0, err
	} else if bytes != 2 {
		return 0, fmt.Errorf("failed to read uint16 (%d instead of 2)", bytes)
	}
	return binary.LittleEndian.Uint16(buf), nil
}

func WriteUint16(i uint16, w io.Writer) error {
	buf := []byte{0, 0}
	binary.LittleEndian.PutUint16(buf, i)
	bytes, err := w.Write(buf)
	if err != nil {
		return err
	} else if bytes != 2 {
		return fmt.Errorf("failed to write uint16 (%d instead of 2)", bytes)
	}
	return nil
}

func ReadUint32(r io.Reader) (uint32, error) {
	buf := []byte{0, 0, 0, 0}
	bytes, err := r.Read(buf)
	if err != nil {
		return 0, err
	} else if bytes != 4 {
		return 0, fmt.Errorf("failed to read uint32 (%d instead of 4)", bytes)
	}
	return binary.LittleEndian.Uint32(buf), nil
}

func WriteUint32(i uint32, w io.Writer) error {
	buf := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(buf, i)
	bytes, err := w.Write(buf)
	if err != nil {
		return err
	} else if bytes != 4 {
		return fmt.Errorf("failed to write uint32 (%d instead of 4)", bytes)
	}
	return nil
}

func ReadUint64(r io.Reader) (uint64, error) {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	bytes, err := r.Read(buf)
	if err != nil {
		return 0, err
	} else if bytes != 8 {
		return 0, fmt.Errorf("failed to read uint32 (%d instead of 8)", bytes)
	}
	return binary.LittleEndian.Uint64(buf), nil
}

func WriteUint64(i uint64, w io.Writer) error {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(buf, i)
	bytes, err := w.Write(buf)
	if err != nil {
		return err
	} else if bytes != 8 {
		return fmt.Errorf("failed to write uint32 (%d instead of 8)", bytes)
	}
	return nil
}

func ReadFloat32(r io.Reader) (float32, error) {
	buf := []byte{0, 0, 0, 0}
	bytes, err := r.Read(buf)
	if err != nil {
		return 0, err
	} else if bytes != 4 {
		return 0, fmt.Errorf("failed to read float32 (%d instead of 4)", bytes)
	}
	bits := binary.LittleEndian.Uint32(buf)
	return math.Float32frombits(bits), nil
}

func WriteFloat32(f float32, w io.Writer) error {
	buf := []byte{0, 0, 0, 0}
	bits := math.Float32bits(f)
	binary.LittleEndian.PutUint32(buf, bits)
	bytes, err := w.Write(buf)
	if err != nil {
		return err
	} else if bytes != 4 {
		return fmt.Errorf("failed to write float32 (%d instead of 4)", bytes)
	}
	return nil
}

func ReadBool(r io.Reader) (bool, error) {
	u, err := ReadUint8(r)
	if u == 0 {
		return false, err
	} else {
		return true, err
	}
}

func WriteBool(b bool, w io.Writer) error {
	if b {
		return WriteUint8(1, w)
	} else {
		return WriteUint8(0, w)
	}
}

func ReadString(r io.Reader) (string, error) {
	size, err := ReadUint32(r)
	if err != nil {
		return "", fmt.Errorf("failed to read string size: %s", err)
	}
	buf := make([]byte, size, size)
	bytes, err := r.Read(buf)
	if err != nil {
		return "", err
	} else if uint32(bytes) != size {
		return "", fmt.Errorf("failed to read string data (%d instead of %d)", bytes, size)
	}
	return string(buf), nil
}

func WriteString(s string, w io.Writer) error {
	if err := WriteUint32(uint32(len(s)), w); err != nil {
		return fmt.Errorf("failed to write string size: %s", err)
	}
	bytes, err := w.Write([]byte(s))
	if err != nil {
		return err
	} else if bytes != len(s) {
		return fmt.Errorf("failed to write string data (%d instead of %d)", bytes, len(s))
	}
	return nil
}
