package wait

import (
	"github.com/stanford-esrg/lzr"
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
	return ""
}


func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "wait", &h )
}

