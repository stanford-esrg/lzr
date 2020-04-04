package ftp

import (
	"lzr"
	"strings"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	data := []byte("") // just wait for banner
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
	if strings.Contains( strings.ToLower(data), "ftp" ) {
		return "ftp"
	}
    return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "ftp", &h )
}

