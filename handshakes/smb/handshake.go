package smb

import (
	"github.com/stanford-esrg/lzr"
	"strings"
	"encoding/hex"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	datas := "00000027ff534d427200000000000000000000000000000000000000000e00024e54204c4d20302e313200"
	data, _:= hex.DecodeString(datas)
	//data := []byte("\x00\x6d\x00\x00\x01\x00\x00\x00\x01\x38\x01\x2c \x00\x6d\x00\x00\x01\x00\x00\x00\x01\x38\x01\x2c")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
	if strings.Contains( data, "SMB" ) {
         return "smb"
	}
    return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "smb", &h )
}

