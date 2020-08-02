package redis

import (
	"github.com/stanford-esrg/lzr"
	"strings"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {

	data := []byte("\x2a\x31\x0d\x0a\x24\x34\x0d\x0a\x50\x49\x4e\x47\x0d\x0a")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
	if (len(data) == 7 && strings.Contains( data, "PONG" )) ||
		strings.Contains( data, "Redis" ) {
         return "redis"
	}
	//thats how it deals with wrong handshake
	if strings.Contains( data, "-ERR unknown" ) {
		return "redis"
	}
    return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "redis", &h )
}

