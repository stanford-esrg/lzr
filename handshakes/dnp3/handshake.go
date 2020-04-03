package dnp3

import (
	"lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}



func (h *HandshakeMod) GetData( dst string ) []byte {
    data := makeLinkRequestBatch(0x0000, 1, 0x0000, 100)
    return data
}

func (h *HandshakeMod) Verify( data string ) string {

	return verifyDNP3( data )

}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "dnp3", &h )
}

