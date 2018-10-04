package util

import (
	"bytes"
	"github.com/lugu/qiloop/type/value"
)

func ErrorPaylad(err error) []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	val := value.String(err.Error())
	val.Write(buf)
	return buf.Bytes()
}
