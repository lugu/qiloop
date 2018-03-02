package basic

import (
    "encoding/binary"
    "io"
    "fmt"
)

func ReadUint32(r io.Reader) (uint32, error) {
    buf := []byte{0, 0, 0, 0}
    bytes, err := r.Read(buf)
    if (err != nil) {
        return 0, err
    } else if (bytes != 4) {
        return 0, fmt.Errorf("failed to read uint32 (%d instead of 4)", bytes)
    }
    return binary.LittleEndian.Uint32(buf), nil
}

func WriteUint32(i uint32, w io.Writer) error {
    buf := []byte{0, 0, 0, 0}
    binary.LittleEndian.PutUint32(buf, i)
    bytes, err := w.Write(buf)
    if (err != nil) {
        return err
    } else if (bytes != 4) {
        return fmt.Errorf("failed to write uint32 (%d instead of 4)", bytes)
    }
    return nil
}

func ReadString(r io.Reader) (string, error) {
    size, err := ReadUint32(r)
    if err != nil {
        return "", fmt.Errorf("failed to read string size: %s", err)
    }
    buf := make([]byte, size, size)
    bytes, err := r.Read(buf)
    if (err != nil) {
        return "", err
    } else if (uint32(bytes) != size) {
        return "", fmt.Errorf("failed to read string data (%d instead of %d)", bytes, size)
    }
    return string(buf), nil
}

func WriteString(s string, w io.Writer) error {
    if err := WriteUint32(uint32(len(s)), w); err != nil {
        return fmt.Errorf("failed to write string size: %s", err)
    }
    bytes, err := w.Write([]byte(s))
    if (err != nil) {
        return err
    } else if (bytes != len(s)) {
        return fmt.Errorf("failed to write string data (%d instead of %d)", bytes, len(s))
    }
    return nil
}

