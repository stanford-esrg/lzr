package telnet

import (
	"lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

//server first protocol
func (h *HandshakeMod) GetData( dst string ) []byte {
	data := []byte("")
	return data
}

func (h *HandshakeMod) Verify( data string ) string {
	if len(data) < 2 {
		return ""
	}
	if ( data[0] != byte(0xff) ) {
		return ""
	}
	if (  data[1] == byte(0xff) || data[1] == byte(0xfe) ||
			data[1] == byte(0xfd) || data[1] == byte(0xfc) ||
				data[1] == byte(0xfb) ) {
		return "telnet"
	}

	return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "telnet", &h )
}

