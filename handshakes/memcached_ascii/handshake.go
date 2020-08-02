package memcached_ascii

import (
	"github.com/stanford-esrg/lzr"
	"strings"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	data := []byte("stats\r\n")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
    if strings.Contains( data, "STAT" ) &&
		strings.Contains( data, "pid" ){
		return "memcached_ascii"
	}
	return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "memcached_ascii", &h )
}

