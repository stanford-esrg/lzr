package amqp

import (
	"github.com/stanford-esrg/lzr"
	"strings"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	data := []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}
    return data
}

func (h *HandshakeMod) Verify( data string ) string {

    if strings.Contains( data, "AMQP" ){
         return "amqp"
    }
    return ""

}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "amqp", &h )
}

