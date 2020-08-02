package siemens

import (
	"github.com/stanford-esrg/lzr"
	"encoding/hex"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	datas := "0300001611e00000000400c1020100c2020200c0010a"
	data, _:= hex.DecodeString(datas)
	//data := []byte("\x00\x6d\x00\x00\x01\x00\x00\x00\x01\x38\x01\x2c \x00\x6d\x00\x00\x01\x00\x00\x00\x01\x38\x01\x2c")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {

    datab := []byte(data)
    if len(datab) < 6 {
        return ""
    }
	size := datab[4]
	pduType := datab[5]
	if int(size) + 1 != len(datab[4:]) || pduType != 0xd0 {
		return ""
	}
	return "siemens"
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "siemens", &h )
}

