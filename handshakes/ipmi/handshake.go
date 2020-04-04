package ipmi

import (
	"lzr"
	//"strings"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {

	data := getIPMI()
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
	return verifyIPMI( data )
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "ipmi", &h )
}

