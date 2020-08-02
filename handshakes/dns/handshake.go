package dns

import (
	"github.com/stanford-esrg/lzr"
	"strings"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	//question: stackoverflow.com
    data := []byte("\x00\x23\x39\x62\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00\x0d\x73\x74\x61\x63\x6b\x6f\x76\x65\x72\x66\x6c\x6f\x77\x03\x63\x6f\x6d\x00\x00\x01\x00\x01")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {

    if strings.Contains( data, "stackoverflow" ){
         return "dns"
    }
    return ""

}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "dns", &h )
}

