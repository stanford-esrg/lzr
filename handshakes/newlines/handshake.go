package newlines

import (
	"github.com/stanford-esrg/lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	data := []byte("\n\n")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
    return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "newlines", &h )
}

