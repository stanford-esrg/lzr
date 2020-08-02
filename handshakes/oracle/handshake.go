package oracle

import (
	"github.com/stanford-esrg/lzr"
	"strings"
	"encoding/hex"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {

	datas := "006d0000010000000138012c0c412000ffff7f08000001000033003a000008004141000000000000000000000000000000000000000000000000284445534352495054494f4e3d28434f4e4e4543545f444154413d284349443d2850524f4752414d3d6c7a724f524129292929"
	data, _:= hex.DecodeString(datas)
	//data := []byte("\x00\x6d\x00\x00\x01\x00\x00\x00\x01\x38\x01\x2c \x00\x6d\x00\x00\x01\x00\x00\x00\x01\x38\x01\x2c")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {
	if strings.Contains( data, "DESCRIPTION=(" ) && strings.Contains( data, "(EMFI=4)" ) {
         return "oracle"
	}
    return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "oracle", &h )
}

