package modbus

import (
	"bytes"
	"github.com/stanford-esrg/lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {

	data := []byte{0x5a,0x47,0x00,0x00,0x00,0x05,0x00,0x2b,0x0e,0x01,0x00}
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
	if len([]byte(data)) < 4 {
		return ""
	}
	if bytes.Equal([]byte(data[0:4]),[]byte{0x5a,0x47,0x00,0x00}) {
		return "modbus"
	}
	return ""

}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "modbus", &h )
}

