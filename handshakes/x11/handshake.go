package x11

import (
	"github.com/stanford-esrg/lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	data := []byte("\x6c\x00\x0b\x00\x00\x00\x00\x00\x00\x00\x00\x00")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
    return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "x11", &h )
}

