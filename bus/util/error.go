package util

import (
	"bytes"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/value"
)

func errorPaylad(err error) []byte {
	var buf bytes.Buffer
	val := value.String(err.Error())
	val.Write(&buf)
	return buf.Bytes()
}

func ReplyError(e net.EndPoint, m *net.Message, err error) error {
	hdr := net.NewHeader(net.Error, m.Header.Service, m.Header.Object,
		m.Header.Action, m.Header.ID)
	mError := net.NewMessage(hdr, errorPaylad(err))
	return e.Send(mError)
}
