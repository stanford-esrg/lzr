package mqtt

import (
	"github.com/stanford-esrg/lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {

	data := []byte{0x10,0xF,0x00,0x04,0x4d,0x51,0x54,0x54,0x04,0x00,0x00,0x0a,0x00,0x03,0x4c,0x5a,0x52}
    return data
}

func (h *HandshakeMod) Verify( data string ) string {

	datab := []byte(data)
	if len(datab) != 4 {
		return ""
	}
	// CTRL_CONNACK
	if byte(datab[0]) == byte(0x20)  {
		//3rd byte ACCEPTED(0x00) or 5 types of refurst(0x01-0x05)
		if int(datab[3]) >= 0x00 && int(datab[3]) <= 0x05 {
			return "mqtt"
		}
	}
	return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "mqtt", &h )
}

