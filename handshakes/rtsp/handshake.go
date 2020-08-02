package rtsp

import (
	"github.com/stanford-esrg/lzr"
	"strings"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	data := []byte("OPTIONS / RTSP/1.0\r\n\r\n")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {

    if strings.Contains( data, "RTSP" ){
         return "rtsp"
    }
    return ""

}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "rtsp", &h )
}

