package pptp

import (
	"github.com/stanford-esrg/lzr"
    "encoding/binary"
	"strings"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {

	msglen := 28+64+64
    message_type := 1 //control message
    magic_cookie := 439041101 //0x 1A2B3C4D
    control_message := 1 // Start-Control-Connection-Request
    version := 256 //version 1
    frame := 1 //Asynchronous
    bearer := 1 //analog
    max_chan := 65535
    firm_rev := 1
    //host_name := ""
    pad := 0
    out := make([]byte, msglen)
    //make msg header
    binary.BigEndian.PutUint16(out[0:],uint16(msglen))
    binary.BigEndian.PutUint16(out[2:],uint16(message_type))
    binary.BigEndian.PutUint32(out[4:],uint32(magic_cookie))
    binary.BigEndian.PutUint16(out[8:],uint16(control_message))
    binary.BigEndian.PutUint16(out[10:],uint16(pad))
    binary.BigEndian.PutUint16(out[12:],uint16(version))
    binary.BigEndian.PutUint16(out[14:],uint16(pad))
    binary.BigEndian.PutUint32(out[16:],uint32(frame))
    binary.BigEndian.PutUint32(out[20:],uint32(bearer))
    binary.BigEndian.PutUint16(out[24:],uint16(max_chan))
    binary.BigEndian.PutUint16(out[26:],uint16(firm_rev))
    binary.BigEndian.PutUint64(out[28:],uint64(pad))
    binary.BigEndian.PutUint64(out[92:],uint64(pad))

    return out

}

func (h *HandshakeMod) Verify( data string ) string {
    if strings.Contains( data, "+<M" ){
         return "pptp"
	}
    return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "pptp", &h )
}

