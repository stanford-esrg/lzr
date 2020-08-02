package mysql

import (
	"bytes"
	"github.com/stanford-esrg/lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	//data := []byte("\x20\x00\x00\x01\x00\x08\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
	data := []byte("")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {


	// standard 4 byte header 5100 0000
    if len([]byte(data)) < 49 {
        return ""
    }
	// sequence number should be 0. skip length 
    if bytes.Equal([]byte(data[3:4]),[]byte{0x00}) && bytes.Equal([]byte(data[4:5]),[]byte{0x0a}){
        return "mysql"
    }
    return ""
}


func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "mysql", &h )
}

