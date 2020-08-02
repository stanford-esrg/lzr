package rdp

import (
	"github.com/stanford-esrg/lzr"
	"bytes"
	"encoding/hex"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {

	datas := "0300002621e00000feca00436f6f6b69653a206d737473686173683d0d0a0100080001000000"

	data, _:= hex.DecodeString(datas)
    return data
}

func (h *HandshakeMod) Verify( data string ) string {

    datab := []byte(data)
	res := []byte("\x03\x00\x00\x13\x0e\xd0\xfe\xca\x12\x34")
    if len(datab) < 11 {
        return ""
    }
	// https://medium.com/@bromiley/what-happens-before-hello-ce9f29fa0cef
	//   [0-3]    TPKT Header
	//   [4-10]   X.224 Class 0 Connection Confirm
	if bytes.Equal([]byte(data[0:10]), res) {
		return "rdp"
    }
    return ""

}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "rdp", &h )
}

