package postgres

import (
	"github.com/stanford-esrg/lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	data := []byte("\x00\x00\x00\x08\x04\xd2\x16\x2f")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
	datab := []byte(data)
	if len(datab) != 1{
		return ""
	}
	// N or S or E
	if byte(datab[0]) == byte(0x4e)  || byte(datab[0]) == byte(0x53) ||
		  byte(datab[0]) == byte(0x45) {
         return "postgres"
	}
    return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "postgres", &h )
}

