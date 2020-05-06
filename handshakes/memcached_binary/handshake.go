package memcached_binary

import (
	"lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	//binary memcached
	data := []byte{0x80, 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
        0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	//data := []byte("stats\r\n")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
	if int([]byte(data)[0]) == 0x81 {
		return "memcached_binary"
	}
	return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "memcached_binary", &h )
}

