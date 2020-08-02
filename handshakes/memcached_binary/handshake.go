package memcached_binary

import (
	"github.com/stanford-esrg/lzr"
	"strings"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	//binary memcached
	data := []byte{0x80, 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
        0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
	if len([]byte(data)) < 1 {
		return ""
	}
	if int([]byte(data)[0]) == 0x81 {
		return "memcached_binary"
	}
	if strings.Contains( data, "ERROR\r\n") {
		return "memcached_binary"
	}
	return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "memcached_binary", &h )
}

