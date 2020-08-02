package mssql

import (
	"github.com/stanford-esrg/lzr"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

func (h *HandshakeMod) GetData( dst string ) []byte {
	//https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-tds/9420b4a3-eb9f-4f5e-90bd-3160444aa5a7
	data := []byte("\x12\x01\x00\x2f\x00\x00\x02\x00\x00\x00\x1a\x00\x06\x01\x00\x02\x00\x01\x02\x00\x21\x00\x01\x03\x00\x22\x00\x04\x04\x00\x26\x00\x01\xff\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00")
    return data
}

func (h *HandshakeMod) Verify( data string ) string {

	datab := []byte(data)
	//https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-tds/517a62a2-7448-47b6-81eb-c0c5027826ca
	//just checking for type and status (04/ 01) b/c the length and spid are hard to validate
	if len(datab) < 6 {
		return ""
	}
	if datab[0] == 0x04 && datab[1] == 0x01 {
			return "mssql"
	}
	return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	lzr.AddHandshake( "mssql", &h )
}

