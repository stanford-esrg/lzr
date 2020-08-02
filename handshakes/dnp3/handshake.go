package dnp3

import (
	"github.com/stanford-esrg/lzr"
    "encoding/binary"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}



func (h *HandshakeMod) GetData( dst string ) []byte {
    return  GetFirstData()
}

func (h *HandshakeMod) Verify( data string ) string {

    data_bytes := []byte(data)
    if len(data_bytes) >= 10 && binary.BigEndian.Uint16(data_bytes[0:2]) == 0x0564 {
        return "dnp3"
    }

    return ""

}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "dnp3", &h )
}

